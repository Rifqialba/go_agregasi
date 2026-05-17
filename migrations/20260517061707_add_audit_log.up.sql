CREATE TABLE audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    action TEXT NOT NULL,

    entity_id TEXT,

    details JSONB,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
REVOKE DELETE, UPDATE ON audit_log FROM app_user;