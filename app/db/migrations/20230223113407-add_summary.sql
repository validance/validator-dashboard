-- +migrate Up
CREATE TABLE delegation_summary
(
    id                                                 SERIAL PRIMARY KEY,
    chain                                              VARCHAR(64)      NOT NULL,

    yesterday_delegation_amount_total                  DOUBLE PRECISION NOT NULL DEFAULT 0,
    yesterday_delegation_amount_b2b                    DOUBLE PRECISION NOT NULL DEFAULT 0,
    yesterday_delegation_amount_b2c                    DOUBLE PRECISION NOT NULL DEFAULT 0,
    yesterday_delegation_amount_unknown                DOUBLE PRECISION NOT NULL DEFAULT 0,

    today_existing_increased_delegation_amount_total   DOUBLE PRECISION NOT NULL DEFAULT 0,
    today_existing_increased_delegation_amount_b2b     DOUBLE PRECISION NOT NULL DEFAULT 0,
    today_existing_increased_delegation_amount_b2c     DOUBLE PRECISION NOT NULL DEFAULT 0,
    today_existing_increased_delegation_amount_unknown DOUBLE PRECISION NOT NULL DEFAULT 0,

    today_new_increased_delegation_amount_total        DOUBLE PRECISION NOT NULL DEFAULT 0,
    today_new_increased_delegation_amount_b2b          DOUBLE PRECISION NOT NULL DEFAULT 0,
    today_new_increased_delegation_amount_b2c          DOUBLE PRECISION NOT NULL DEFAULT 0,
    today_new_increased_delegation_amount_unknown      DOUBLE PRECISION NOT NULL DEFAULT 0,

    today_return_increased_delegation_amount_total     DOUBLE PRECISION NOT NULL DEFAULT 0,
    today_return_increased_delegation_amount_b2b       DOUBLE PRECISION NOT NULL DEFAULT 0,
    today_return_increased_delegation_amount_b2c       DOUBLE PRECISION NOT NULL DEFAULT 0,
    today_return_increased_delegation_amount_unknown   DOUBLE PRECISION NOT NULL DEFAULT 0,

    today_existing_decreased_delegation_amount_total   DOUBLE PRECISION NOT NULL DEFAULT 0,
    today_existing_decreased_delegation_amount_b2b     DOUBLE PRECISION NOT NULL DEFAULT 0,
    today_existing_decreased_delegation_amount_b2c     DOUBLE PRECISION NOT NULL DEFAULT 0,
    today_existing_decreased_delegation_amount_unknown DOUBLE PRECISION NOT NULL DEFAULT 0,

    today_left_decreased_delegation_amount_total       DOUBLE PRECISION NOT NULL DEFAULT 0,
    today_left_decreased_delegation_amount_b2b         DOUBLE PRECISION NOT NULL DEFAULT 0,
    today_left_decreased_delegation_amount_b2c         DOUBLE PRECISION NOT NULL DEFAULT 0,
    today_left_decreased_delegation_amount_unknown     DOUBLE PRECISION NOT NULL DEFAULT 0,

    create_dt                                          TIMESTAMP        NOT NULL DEFAULT current_timestamp
);

CREATE INDEX IF NOT EXISTS idx_delegation_summary_create_dt ON delegation_summary (create_dt);

-- +migrate Down

DROP INDEX IF EXISTS idx_delegation_summary_create_dt;
DROP TABLE IF EXISTS delegation_summary;
