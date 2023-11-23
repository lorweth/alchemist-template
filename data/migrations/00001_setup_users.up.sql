--
-- USERS table
--
CREATE TABLE IF NOT EXISTS "users" (
    "id"            BIGINT PRIMARY KEY,
    "email"         VARCHAR(80) NOT NULL,
    "created_at"    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "deleted_at"    TIMESTAMPTZ NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS "email_on_users" ON "users"("email");
