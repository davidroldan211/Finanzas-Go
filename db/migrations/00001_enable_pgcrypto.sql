-- +goose Up
-- Necesaria para gen_random_uuid() en las PKs.
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- +goose Down
-- No-op a propósito: otras bases del mismo cluster pueden depender de pgcrypto,
-- así que no se desinstala al hacer rollback.
