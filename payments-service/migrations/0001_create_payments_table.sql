CREATE TABLE payments (
    id             UUID PRIMARY KEY,
    merchant_id    UUID NOT NULL,
    amount         BIGINT NOT NULL,
    currency       TEXT NOT NULL,
    status         TEXT NOT NULL,
    created_at     TIMESTAMPTZ NOT NULL,
    history        JSONB NOT NULL
);