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

CREATE TABLE "named_keys"
(
    "uref"          VARCHAR(77) PRIMARY KEY,
    "name"          VARCHAR NOT NULL,
    "is_purse"      BOOLEAN NOT NULL,
    "initial_value" jsonb
);

CREATE TABLE "contracts_named_keys"
(
    contract_hash  VARCHAR(64) references contracts (hash),
    named_key_uref VARCHAR(77) references named_keys (uref),
    primary key (contract_hash, named_key_uref)
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
    "staked_amount"   NUMERIC     NOT NULL,
    "delegation_rate" INT         NOT NULL,
    "inactive"        BOOL        NOT NULL
);

CREATE TABLE "delegators"
(
    "public_key"    VARCHAR(68) NOT NULL,
    "delegatee"     VARCHAR(68) NOT NULL,
    "staked_amount" NUMERIC     NOT NULL,
    "bonding_purse" VARCHAR     NOT NULL
);

CREATE TABLE "accounts"
(
    "account_hash" VARCHAR(64) NOT NULL PRIMARY KEY,
    "public_key"   VARCHAR(68) UNIQUE,
    "main_purse"   VARCHAR(73) NOT NULL UNIQUE
);

CREATE TABLE "purses"
(
    "purse"   VARCHAR(73) NOT NULL PRIMARY KEY,
    "balance" NUMERIC
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
SELECT sum(amount::NUMERIC) as total_rewards
FROM rewards;

CREATE VIEW stakers AS
WITH publicKeys as (SELECT DISTINCT public_key
                    FROM delegators)
SELECT COUNT(*)
from publicKeys;

CREATE VIEW mouvements AS
SELECT 'delegate'                                                as type,
       FLOOR(SUM((metadata ->> 'amount')::numeric) / 1000000000) as count,
       date_trunc('day', timestamp)                              as day
from deploys
WHERE timestamp >= NOW() - INTERVAL '14 DAY'
  and metadata_type = 'delegate'
  AND result is true
GROUP BY day
UNION
SELECT 'undelegate'                                              as type,
       FLOOR(SUM((metadata ->> 'amount')::numeric) / 1000000000) as count,
       date_trunc('day', timestamp)                              as day
from deploys
WHERE timestamp >= NOW() - INTERVAL '14 DAY'
  and metadata_type = 'undelegate'
  AND result is true
GROUP BY day
UNION
SELECT 'transfer'                                                as type,
       FLOOR(SUM((metadata ->> 'amount')::numeric) / 1000000000) as count,
       date_trunc('day', timestamp)                              as day
from deploys
WHERE timestamp >= NOW() - INTERVAL '14 DAY'
  and type = 'transfer'
  AND result is true
GROUP BY day;

CREATE VIEW rich_list AS
SELECT *
from purses
         FULL JOIN accounts ON purses.purse = accounts.main_purse
ORDER BY balance desc;

CREATE FUNCTION era_rewards(eraid integer) RETURNS NUMERIC AS
$$
SELECT sum(amount::NUMERIC)
FROM rewards
where era = eraid;
$$ LANGUAGE SQL;

CREATE FUNCTION total_validator_rewards(publickey VARCHAR(68), OUT validator_rewards NUMERIC,
                                        OUT total_rewards NUMERIC) AS
$$
SELECT sum(amount::NUMERIC)                 as total_rewards,
       (SELECT sum(amount::NUMERIC)
        FROM rewards
        where validator_public_key = publickey
          and delegator_public_key is null) as validator_rewards
FROM rewards
where validator_public_key = publickey;
$$ LANGUAGE SQL;

CREATE FUNCTION total_account_rewards(publickey VARCHAR(68)) RETURNS NUMERIC AS
$$
SELECT sum(amount::NUMERIC)
FROM rewards
where delegator_public_key = publickey;
$$ LANGUAGE SQL;

CREATE FUNCTION block_details(blockhash VARCHAR(64), OUT total NUMERIC, OUT success NUMERIC, OUT failed NUMERIC,
                              OUT total_cost NUMERIC) AS
$$
SELECT count(*)                                                                   as total,
       (SELECT count(*) from deploys where block = blockhash and result is true)  as success,
       (SELECT count(*) from deploys where block = blockhash and result is false) as failed,
       sum(cost::NUMERIC)                                                         as total_cost
FROM deploys
where block = blockhash;
$$ LANGUAGE SQL;

CREATE FUNCTION contract_details(contracthash VARCHAR(64), OUT total NUMERIC, OUT success NUMERIC, OUT failed NUMERIC,
                                 OUT total_cost NUMERIC) AS
$$
SELECT count(*)                                                                              as total,
       (SELECT count(*) from deploys where contract_hash = contracthash and result is true)  as success,
       (SELECT count(*) from deploys where contract_hash = contracthash and result is false) as failed,
       sum(cost::NUMERIC)                                                                    as total_cost
FROM deploys
where contract_hash = contracthash;
$$ LANGUAGE SQL;

CREATE ROLE web_anon NOLOGIN;

grant usage on schema public to web_anon;
grant select on public.blocks to web_anon;
grant select on public.raw_blocks to web_anon;
grant select on public.deploys to web_anon;
grant select on public.raw_deploys to web_anon;
grant select on public.contract_packages to web_anon;
grant select on public.contracts to web_anon;
grant select on public.named_keys to web_anon;
grant select on public.contracts_named_keys to web_anon;
grant select on public.rewards to web_anon;
grant select on public.bids to web_anon;
grant select on public.delegators to web_anon;
grant select on public.accounts to web_anon;
grant select on public.purses to web_anon;
grant select on public.full_stats to web_anon;
grant select on public.simple_stats to web_anon;
grant select on public.total_rewards to web_anon;
grant select on public.stakers to web_anon;
grant select on public.mouvements to web_anon;
grant select on public.rich_list to web_anon;
grant execute on function era_rewards(integer) to web_anon;
grant execute on function total_validator_rewards(VARCHAR(68)) to web_anon;
grant execute on function total_account_rewards(VARCHAR(68)) to web_anon;
grant execute on function block_details(VARCHAR(64)) to web_anon;
grant execute on function contract_details(VARCHAR(64)) to web_anon;
