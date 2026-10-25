-- Copyright 2026 The ThunderID Authors
-- SPDX-License-Identifier: Apache-2.0

-- ============================================================
-- Stored procedure: purge expired runtime_transient rows in bounded batches.
--
-- Deletes expired rows in batches of p_batch_size (default 1000), committing
-- after each batch to keep locks short on large tables. Must run as a top-level
-- CALL (the per-batch COMMIT cannot run inside an outer transaction).
--
-- Run once manually (ad-hoc / on-demand):
--   PGPASSWORD=<pass> psql -h <host> -p <port> -U <user> -d <runtime_transient> \
--     -c "CALL cleanup_expired_runtime_transient_data();"
--
--   -- Optional: override the batch size (rows deleted per batch):
--   -c "CALL cleanup_expired_runtime_transient_data(500);"
--
-- Scheduled execution options:
--
--   1. pg_cron (RECOMMENDED, requires the pg_cron extension):
--      CREATE EXTENSION IF NOT EXISTS pg_cron;
--      SELECT cron.schedule(
--        'cleanup-runtime_transient-expired',
--        '*/60 * * * *',
--        $$CALL cleanup_expired_runtime_transient_data()$$
--      );
--      -- To verify: SELECT * FROM cron.job WHERE jobname = 'cleanup-runtime_transient-expired';
--      -- To remove: SELECT cron.unschedule('cleanup-runtime_transient-expired');
--
--   2. Kubernetes CronJob: call CALL cleanup_expired_runtime_transient_data()
--      via a psql container on the desired schedule.
--
--   3. OS cron (every 60 minutes):
-- --      */60 * * * * postgres PGPASSWORD=<pass> psql -h <host> -p <port> \
-- --        -U <user> -d <runtime_transient> -c "CALL cleanup_expired_runtime_transient_data();" \
-- --        >> /var/log/thunderid-cleanup.log 2>&1
-- ============================================================

-- Drop the old parameterless signature so re-applying doesn't leave an ambiguous overload.
DROP PROCEDURE IF EXISTS cleanup_expired_runtime_transient_data();

CREATE OR REPLACE PROCEDURE cleanup_expired_runtime_transient_data(p_batch_size INT DEFAULT 1000)
LANGUAGE plpgsql
AS $$
DECLARE
    v_now     TIMESTAMP := NOW() AT TIME ZONE 'UTC';
    v_deleted INT;
BEGIN
    -- Guard against a batch size that would disable batching or make no progress.
    IF p_batch_size IS NULL OR p_batch_size <= 0 THEN
        p_batch_size := 1000;
    END IF;

    -- RUNTIME_STORE is LIST-partitioned by NAMESPACE, where ctid is not unique across
    -- partitions; match rows by primary key. ORDER BY EXPIRY_TIME lets the batch be
    -- located via an index scan.
    LOOP
        DELETE FROM "RUNTIME_STORE"
        WHERE (DEPLOYMENT_ID, NAMESPACE, KEY) IN (
            SELECT DEPLOYMENT_ID, NAMESPACE, KEY FROM "RUNTIME_STORE"
            WHERE EXPIRY_TIME < v_now ORDER BY EXPIRY_TIME LIMIT p_batch_size
        );
        GET DIAGNOSTICS v_deleted = ROW_COUNT;
        COMMIT;
        EXIT WHEN v_deleted = 0;
    END LOOP;
END;
$$;
