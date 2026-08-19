-- Modify "users" table
ALTER TABLE "public"."users" DROP CONSTRAINT "uni_users_pin", ALTER COLUMN "pin" DROP NOT NULL;
-- Create index "idx_users_phone_number" to table: "users"
CREATE UNIQUE INDEX "idx_users_phone_number" ON "public"."users" ("phone_number");
-- Create index "idx_users_pin" to table: "users"
CREATE UNIQUE INDEX "idx_users_pin" ON "public"."users" ("pin");
-- Create index "idx_users_science_id" to table: "users"
CREATE UNIQUE INDEX "idx_users_science_id" ON "public"."users" ("science_id");
-- Create "otp_codes" table
CREATE TABLE "public"."otp_codes" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "phone" character varying(50) NOT NULL,
  "code" character varying(32) NOT NULL,
  "session_id" character varying(255) NOT NULL,
  "purpose" character varying(32) NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "used_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_otp_codes_deleted_at" to table: "otp_codes"
CREATE INDEX "idx_otp_codes_deleted_at" ON "public"."otp_codes" ("deleted_at");
-- Drop "telegram_auth_sessions" table
DROP TABLE "public"."telegram_auth_sessions";
