-- Table to store generic runtime key-value entries, isolated by NAMESPACE.
CREATE TABLE "RUNTIME_STORE" (
    DEPLOYMENT_ID VARCHAR(255) NOT NULL,
    NAMESPACE     VARCHAR(64)  NOT NULL,
    KEY           VARCHAR(512) NOT NULL,
    VALUE         TEXT         NOT NULL,
    EXPIRY_TIME   DATETIME,
    CREATED_AT    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UPDATED_AT    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (DEPLOYMENT_ID, NAMESPACE, KEY)
);

-- Index for expiry time on RUNTIME_STORE (supports cleanup and expiry checks)
CREATE INDEX idx_runtime_store_expiry_time ON "RUNTIME_STORE" (EXPIRY_TIME);
