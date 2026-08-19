-- Create "clients" table
CREATE TABLE "public"."clients" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "client_id" character varying(255) NOT NULL,
  "client_secret" character varying(255) NOT NULL,
  "name" character varying(255) NOT NULL,
  "description" text NULL,
  "callback_url" text NULL,
  "is_active" boolean NULL DEFAULT true,
  PRIMARY KEY ("id")
);
-- Create index "idx_clients_client_id" to table: "clients"
CREATE UNIQUE INDEX "idx_clients_client_id" ON "public"."clients" ("client_id");
-- Create index "idx_clients_deleted_at" to table: "clients"
CREATE INDEX "idx_clients_deleted_at" ON "public"."clients" ("deleted_at");
