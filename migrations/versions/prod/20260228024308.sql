-- Modify "users" table
ALTER TABLE "public"."users" ALTER COLUMN "orc_id_id" DROP DEFAULT;
-- Create index "idx_users_orc_id_id" to table: "users"
CREATE UNIQUE INDEX "idx_users_orc_id_id" ON "public"."users" ("orc_id_id");
-- Create "telegram_auth_sessions" table
CREATE TABLE "public"."telegram_auth_sessions" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "session_id" uuid NULL DEFAULT gen_random_uuid(),
  "science_id" character varying(50) NULL,
  "verified_phone" character varying(20) NULL,
  "telegram_chat_id" bigint NULL,
  "user_id" bigint NULL,
  "status" bigint NULL DEFAULT 10,
  "expires_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_telegram_auth_sessions_deleted_at" to table: "telegram_auth_sessions"
CREATE INDEX "idx_telegram_auth_sessions_deleted_at" ON "public"."telegram_auth_sessions" ("deleted_at");
-- Create index "idx_telegram_auth_sessions_expires_at" to table: "telegram_auth_sessions"
CREATE INDEX "idx_telegram_auth_sessions_expires_at" ON "public"."telegram_auth_sessions" ("expires_at");
-- Create index "idx_telegram_auth_sessions_science_id" to table: "telegram_auth_sessions"
CREATE INDEX "idx_telegram_auth_sessions_science_id" ON "public"."telegram_auth_sessions" ("science_id");
-- Create index "idx_telegram_auth_sessions_telegram_chat_id" to table: "telegram_auth_sessions"
CREATE INDEX "idx_telegram_auth_sessions_telegram_chat_id" ON "public"."telegram_auth_sessions" ("telegram_chat_id");
-- Create index "idx_telegram_auth_sessions_user_id" to table: "telegram_auth_sessions"
CREATE INDEX "idx_telegram_auth_sessions_user_id" ON "public"."telegram_auth_sessions" ("user_id");
-- Create "journal_configs" table
CREATE TABLE "public"."journal_configs" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "journal_id" bigint NULL,
  "website_url" text NULL,
  "conf" jsonb NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_journal_configs_journal" FOREIGN KEY ("journal_id") REFERENCES "public"."journals" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_journal_configs_deleted_at" to table: "journal_configs"
CREATE INDEX "idx_journal_configs_deleted_at" ON "public"."journal_configs" ("deleted_at");
-- Create index "idx_journal_configs_journal_id" to table: "journal_configs"
CREATE UNIQUE INDEX "idx_journal_configs_journal_id" ON "public"."journal_configs" ("journal_id");
-- Create index "idx_journal_configs_website_url" to table: "journal_configs"
CREATE UNIQUE INDEX "idx_journal_configs_website_url" ON "public"."journal_configs" ("website_url");
