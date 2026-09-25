-- dev/main: 175_user_totp owns totp_secret_encrypted / totp_enabled_at here,
-- so rolling back 551 only removes the columns it introduced.
ALTER TABLE "user" DROP COLUMN IF EXISTS totp_locked_until;
ALTER TABLE "user" DROP COLUMN IF EXISTS totp_failed_attempts;
ALTER TABLE "user" DROP COLUMN IF EXISTS totp_last_used_step;
