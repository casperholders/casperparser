CREATE TABLE "blocks"
(
    "hash"      VARCHAR(64) PRIMARY KEY,
    "era"       BIGINT      NOT NULL,
    "timestamp" timestamptz NOT NULL,
    "height"    BIGINT      NOT NULL,
    "era_end"   bool        NOT NULL,
    "validated" bool        NOT NULL
);

CREATE TABLE "raw_blocks"
(
    "hash" VARCHAR(64) PRIMARY KEY,
    "data" jsonb NOT NULL
);

CREATE TABLE "deploys"
(
    "hash"          VARCHAR(64) PRIMARY KEY,
    "from"          VARCHAR(68) NOT NULL,
    "cost"          VARCHAR     NOT NULL,
    "result"        boolean     NOT NULL,
    "timestamp"     timestamptz NOT NULL,
    "block"         VARCHAR(64) NOT NULL,
    "type"          VARCHAR     NOT NULL,
    "metadata_type" VARCHAR     NOT NULL,
    "contract_hash" VARCHAR(64),
    "contract_name" VARCHAR,
    "entrypoint"    VARCHAR,
    "metadata"      jsonb,
    "events"        jsonb
);

CREATE TABLE "raw_deploys"
(
    "hash" VARCHAR(64) PRIMARY KEY,
    "data" jsonb NOT NULL
);

CREATE TABLE "contract_packages"
(
    "hash"   VARCHAR(64) PRIMARY KEY,
    "deploy" VARCHAR(64),
    "from"   VARCHAR(68),
    "data"   jsonb NOT NULL
);

CREATE TABLE "contracts"
(
    "hash"    VARCHAR(64) PRIMARY KEY,
    "package" VARCHAR(64) NOT NULL,
    "deploy"  VARCHAR(64),
    "from"    VARCHAR(68),
    "type"    VARCHAR NOT NULL,
    "data"    jsonb   NOT NULL
);

ALTER TABLE "deploys"
    ADD FOREIGN KEY ("block") REFERENCES "blocks" ("hash");

ALTER TABLE "blocks"
    ADD FOREIGN KEY ("hash") REFERENCES "raw_blocks" ("hash");

ALTER TABLE "deploys"
    ADD FOREIGN KEY ("hash") REFERENCES "raw_deploys" ("hash");

ALTER TABLE "contracts"
    ADD FOREIGN KEY ("package") REFERENCES "contract_packages" ("hash");

CREATE INDEX ON "deploys" ("block");
CREATE INDEX ON "deploys" ("from");

CREATE VIEW full_stats AS
SELECT count(*), type, date_trunc('day', timestamp) as day
from deploys
WHERE timestamp >= NOW() - INTERVAL '14 DAY'
GROUP BY day, type;
