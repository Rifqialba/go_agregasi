CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE raw_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    source_id TEXT NOT NULL,

    source_type TEXT NOT NULL,

    raw_payload JSONB NOT NULL,

    received_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    status TEXT NOT NULL DEFAULT 'PENDING',

    idempotency_key TEXT UNIQUE
);