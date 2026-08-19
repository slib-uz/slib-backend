-- Create "sandbox_users" table
CREATE TABLE "public"."sandbox_users" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "full_name" text NULL,
  "science_id" text NULL,
  "phone_number" text NULL,
  "otp" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_sandbox_users_deleted_at" to table: "sandbox_users"
CREATE INDEX "idx_sandbox_users_deleted_at" ON "public"."sandbox_users" ("deleted_at");
