-- Adds per-user TOTP state. IF NOT EXISTS keeps this compatible with personal
-- deployments that previously applied the feature under 120_user_totp.
-- NULL totp_secret_encrypted = not set up. A non-NULL secret with a NULL
-- enabled timestamp means setup is waiting for its first verification.
ALTER TABLE "user" ADD COLUMN IF NOT EXISTS totp_secret_encrypted BYTEA;
ALTER TABLE "user" ADD COLUMN IF NOT EXISTS totp_enabled_at TIMESTAMPTZ;
