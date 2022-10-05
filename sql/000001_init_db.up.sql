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
    "type"    VARCHAR     NOT NULL,
    "score"   FLOAT       NOT NULL,
    "data"    jsonb       NOT NULL
);


CREATE TABLE "rewards"
(
    "block"                VARCHAR(64) NOT NULL,
    "era"                  BIGINT      NOT NULL,
    "delegator_public_key" VARCHAR(68),
    "validator_public_key" VARCHAR(68) NOT NULL,
    "amount"               VARCHAR     NOT NULL
);

CREATE TABLE "bids"
(
    "public_key"      VARCHAR(68) NOT NULL PRIMARY KEY,
    "bonding_purse"   VARCHAR     NOT NULL,
    "staked_amount"   VARCHAR     NOT NULL,
    "delegation_rate" INT         NOT NULL,
    "inactive"        BOOL        NOT NULL
);

CREATE TABLE "delegators"
(
    "public_key"    VARCHAR(68) NOT NULL,
    "delegatee"     VARCHAR(68) NOT NULL,
    "staked_amount" VARCHAR     NOT NULL,
    "bonding_purse" VARCHAR     NOT NULL
);

ALTER TABLE "delegators"
    ADD CONSTRAINT uAuction UNIQUE (public_key, delegatee, bonding_purse);

ALTER TABLE "rewards"
    ADD FOREIGN KEY ("block") REFERENCES "blocks" ("hash");

ALTER TABLE "rewards"
    ADD CONSTRAINT uReward UNIQUE (block, era, delegator_public_key, validator_public_key);

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
CREATE INDEX ON "deploys" ("contract_hash");
CREATE INDEX ON "deploys" ("result");
CREATE INDEX ON "deploys" ("timestamp");
CREATE INDEX ON "delegators" ("delegatee");

CREATE VIEW full_stats AS
SELECT count(*), type, date_trunc('day', timestamp) as day
from deploys
WHERE timestamp >= NOW() - INTERVAL '14 DAY'
GROUP BY day, type;

CREATE VIEW simple_stats AS
SELECT count(*), date_trunc('day', timestamp) as day
from deploys
WHERE timestamp >= NOW() - INTERVAL '14 DAY'
GROUP BY day;

CREATE VIEW total_rewards AS
SELECT sum(amount::BIGINT) as total_rewards
FROM rewards;

CREATE FUNCTION era_rewards(eraid integer) RETURNS BIGINT AS
$$
SELECT sum(amount::BIGINT)
FROM rewards
where era = eraid;
$$ LANGUAGE SQL;

CREATE FUNCTION total_validator_rewards(publickey VARCHAR(66), OUT validator_rewards bigint,
                                        OUT total_rewards bigint) AS
$$
SELECT sum(amount::BIGINT)                  as total_rewards,
       (SELECT sum(amount::BIGINT)
        FROM rewards
        where validator_public_key = publickey
          and delegator_public_key is null) as validator_rewards
FROM rewards
where validator_public_key = publickey;
$$ LANGUAGE SQL;

CREATE FUNCTION total_account_rewards(publickey VARCHAR(66)) RETURNS BIGINT AS
$$
SELECT sum(amount::BIGINT)
FROM rewards
where delegator_public_key = publickey;
$$ LANGUAGE SQL;

CREATE FUNCTION block_details(blockhash VARCHAR(64), OUT total bigint, OUT success bigint, OUT failed bigint,
                              OUT total_cost bigint) AS
$$
SELECT count(*)                                                                   as total,
       (SELECT count(*) from deploys where block = blockhash and result is true)  as success,
       (SELECT count(*) from deploys where block = blockhash and result is false) as failed,
       sum(cost::BIGINT)                                                          as total_cost
FROM deploys
where block = blockhash;
$$ LANGUAGE SQL;

CREATE FUNCTION contract_details(contracthash VARCHAR(64), OUT total bigint, OUT success bigint, OUT failed bigint,
                                 OUT total_cost bigint) AS
$$
SELECT count(*)                                                                              as total,
       (SELECT count(*) from deploys where contract_hash = contracthash and result is true)  as success,
       (SELECT count(*) from deploys where contract_hash = contracthash and result is false) as failed,
       sum(cost::BIGINT)                                                                     as total_cost
FROM deploys
where contract_hash = contracthash;
$$ LANGUAGE SQL;

CREATE ROLE web_anon NOLOGIN;

grant usage on schema public to web_anon;
grant select on public.blocks to web_anon;
grant select on public.contract_packages to web_anon;
grant select on public.contracts to web_anon;
grant select on public.deploys to web_anon;
grant select on public.raw_blocks to web_anon;
grant select on public.raw_deploys to web_anon;
grant select on public.full_stats to web_anon;
grant select on public.simple_stats to web_anon;
grant select on public.rewards to web_anon;
grant select on public.total_rewards to web_anon;
grant select on public.bids to web_anon;
grant select on public.delegators to web_anon;
grant execute on function era_rewards(integer) to web_anon;
grant execute on function block_details(VARCHAR(64)) to web_anon;
grant execute on function contract_details(VARCHAR(64)) to web_anon;
grant execute on function total_validator_rewards(VARCHAR(66)) to web_anon;
grant execute on function total_account_rewards(VARCHAR(66)) to web_anon;