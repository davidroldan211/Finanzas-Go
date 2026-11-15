-- +goose Up
-- NULL = correo sin confirmar. El registro (SPEC 01) la deja en NULL; la
-- confirmación del código (HU-02) es quien la llena.
ALTER TABLE users ADD COLUMN email_verified_at TIMESTAMPTZ NULL;

-- +goose Down
ALTER TABLE users DROP COLUMN IF EXISTS email_verified_at;
