BEGIN;
DO $pair$
DECLARE parent_id text; child_id text; n integer;
BEGIN
 SELECT COUNT(*) INTO n FROM esignet.client_detail WHERE id LIKE 'setupass-native-parent-%' AND status='ACTIVE';
 IF n<>1 THEN RAISE EXCEPTION 'Expected one native parent. Open Setu ID and wait for Ready. Found %',n; END IF;
 SELECT COUNT(*) INTO n FROM esignet.client_detail WHERE id LIKE 'setupass-native-child-%' AND status='ACTIVE';
 IF n<>1 THEN RAISE EXCEPTION 'Expected one native child. Open Setu Benefits and wait for Ready. Found %',n; END IF;
 SELECT id INTO parent_id FROM esignet.client_detail WHERE id LIKE 'setupass-native-parent-%' AND status='ACTIVE';
 SELECT id INTO child_id FROM esignet.client_detail WHERE id LIKE 'setupass-native-child-%' AND status='ACTIVE';
 UPDATE esignet.client_detail SET additional_config=(COALESCE(NULLIF(additional_config,''),'{}')::jsonb || jsonb_build_object('native_sso_family','setupass-local-approved'))::text WHERE id IN (parent_id,child_id);
 RAISE NOTICE 'Approved Setu ID and Setu Benefits for local native SSO';
END $pair$;
COMMIT;
