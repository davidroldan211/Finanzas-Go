CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE email_verifications (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email          TEXT        NOT NULL,
    code_hash      TEXT        NOT NULL,
    expires_at     TIMESTAMPTZ NOT NULL,
    attempts       INT         NOT NULL DEFAULT 0,
    max_attempts   INT         NOT NULL DEFAULT 5,
    cooldown_until TIMESTAMPTZ NULL,
    used_at        TIMESTAMPTZ NULL,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_email_verifications_email ON email_verifications (email);
