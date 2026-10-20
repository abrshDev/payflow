CREATE TABLE ledger_entries (
    id               UUID PRIMARY KEY,
    payment_id       UUID NOT NULL,
    amount           BIGINT NOT NULL,
    currency         TEXT NOT NULL,
    direction        TEXT NOT NULL,
    source_event_id  TEXT NOT NULL UNIQUE,
    created_at       TIMESTAMPTZ NOT NULL
);