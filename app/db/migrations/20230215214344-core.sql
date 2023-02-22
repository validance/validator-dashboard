-- +migrate Up
CREATE TYPE address_label AS ENUM ('a41', 'a41ventures', 'grant', 'b2b', 'b2c', 'unknown');
CREATE TYPE address_type AS ENUM ('new', 'existing', 'leave', 'return');

CREATE TABLE IF NOT EXISTS delegation_history
(
    id        SERIAL PRIMARY KEY,
    address   VARCHAR(256) NOT NULL,
    validator VARCHAR(256),
    chain     VARCHAR(64)  NOT NULL,
    amount    DOUBLE PRECISION NOT NULL,
    create_dt TIMESTAMP    NOT NULL default current_timestamp
);

CREATE INDEX IF NOT EXISTS idx_delegation_history_address ON delegation_history (address);
CREATE INDEX IF NOT EXISTS idx_delegation_history_create_dt ON delegation_history (create_dt);
CREATE INDEX IF NOT EXISTS idx_delegation_history_chain_create_dt ON delegation_history (chain, create_dt);
CREATE INDEX IF NOT EXISTS idx_delegation_history_address_chain_create_dt ON delegation_history (address, chain, create_dt);

CREATE TABLE IF NOT EXISTS address_status
(
    id         SERIAL PRIMARY KEY,
    address    VARCHAR(256)  NOT NULL,
    chain      VARCHAR(64)   NOT NULL,
    label      address_label NOT NULL DEFAULT 'unknown',
    status     address_type  NOT NULL DEFAULT 'new',
    create_dt  TIMESTAMP     NOT NULL DEFAULT current_timestamp,
    update_dt TIMESTAMP     NOT NULL DEFAULT current_timestamp
);

CREATE INDEX IF NOT EXISTS idx_address_status_address ON address_status (address);
CREATE INDEX IF NOT EXISTS idx_address_status_create_dt ON address_status (create_dt);
CREATE INDEX IF NOT EXISTS idx_address_status_update_dt ON address_status (update_dt);

CREATE TABLE IF NOT EXISTS income_history
(
    id         SERIAL PRIMARY KEY,
    address    VARCHAR(256) NOT NULL,
    chain      VARCHAR(64)  NOT NULL,
    reward     VARCHAR(256) NOT NULL,
    commission VARCHAR(256) NOT NULL,
    create_dt  TIMESTAMP    NOT NULL DEFAULT current_timestamp
);

CREATE INDEX idx_income_history_create_dt ON income_history (create_dt);
CREATE INDEX idx_income_history_address_chain_create_dt ON income_history (address, chain, create_dt);

CREATE TABLE IF NOT EXISTS token_price
(
    id        SERIAL PRIMARY KEY,
    chain     VARCHAR(64) NOT NULL,
    ticker    VARCHAR(64) NOT NULL,
    price     FLOAT       NOT NULL,
    create_dt TIMESTAMP   NOT NULL DEFAULT current_timestamp
);

CREATE INDEX idx_token_price_created_dt ON token_price (create_dt);
CREATE INDEX idx_token_price_chain_price_created_dt ON token_price (chain, price, create_dt);
CREATE INDEX idx_token_price_ticker_price_created_dt ON token_price (chain, ticker, create_dt);

CREATE TABLE IF NOT EXISTS grant_reward_history
(
    id            SERIAL PRIMARY KEY,
    grant_address VARCHAR(256) NOT NULL,
    validator     VARCHAR(256) NOT NULL,
    chain         VARCHAR(64)  NOT NULL,
    reward        VARCHAR(256) NOT NULL,
    create_dt     TIMESTAMP    NOT NULL DEFAULT current_timestamp
);

CREATE INDEX idx_grant_reward_history_create_dt ON grant_reward_history (create_dt);
CREATE INDEX idx_grant_reward_history_grant_address_create_dt ON grant_reward_history (grant_address, create_dt);

-- +migrate Down
DROP INDEX IF EXISTS idx_grant_reward_history_create_dt;
DROP INDEX IF EXISTS idx_grant_reward_history_grant_address_create_dt;
DROP TABLE IF EXISTS grant_reward_history;

DROP INDEX IF EXISTS idx_token_price_created_dt;
DROP INDEX IF EXISTS idx_token_price_chain_price_created_dt;
DROP INDEX IF EXISTS idx_token_price_ticker_price_created_dt;
DROP TABLE IF EXISTS token_price;

DROP INDEX IF EXISTS idx_income_history_create_dt;
DROP INDEX IF EXISTS idx_income_history_address_chain_create_dt;
DROP TABLE IF EXISTS income_history;

DROP INDEX IF EXISTS idx_address_status_create_dt;
DROP INDEX IF EXISTS idx_address_status_update_dt;
DROP TABLE IF EXISTS address_status;

DROP INDEX IF EXISTS idx_delegation_history_address;
DROP INDEX IF EXISTS idx_delegation_history_create_dt;
DROP INDEX IF EXISTS idx_delegation_history_chain_create_dt;
DROP INDEX IF EXISTS idx_delegation_history_address_chain_create_dt;
DROP TABLE IF EXISTS delegation_history;

DROP TYPE IF EXISTS address_label;
DROP TYPE IF EXISTS address_type;