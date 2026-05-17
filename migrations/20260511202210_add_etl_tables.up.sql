CREATE TABLE processed_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    raw_data_id UUID NOT NULL,

    source_id TEXT NOT NULL,

    normalized_payload JSONB NOT NULL,

    processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    idempotency_key TEXT UNIQUE NOT NULL
);

CREATE TABLE validation_errors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    raw_data_id UUID NOT NULL,

    source_id TEXT NOT NULL,

    error_message TEXT NOT NULL,

    raw_payload JSONB NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);