-- Create "doi_mails" table
CREATE TABLE "public"."doi_mails" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "message_id" text NOT NULL,
  "subject" text NOT NULL,
  "body" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_doi_mails_deleted_at" to table: "doi_mails"
CREATE INDEX "idx_doi_mails_deleted_at" ON "public"."doi_mails" ("deleted_at");
-- Create index "idx_doi_mails_message_id" to table: "doi_mails"
CREATE UNIQUE INDEX "idx_doi_mails_message_id" ON "public"."doi_mails" ("message_id");
-- Create "sequences" table
CREATE TABLE "public"."sequences" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "key" text NOT NULL,
  "seq" bigint NOT NULL,
  "idempotency_key" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_sequences_deleted_at" to table: "sequences"
CREATE INDEX "idx_sequences_deleted_at" ON "public"."sequences" ("deleted_at");
-- Create index "idx_sequences_key" to table: "sequences"
CREATE UNIQUE INDEX "idx_sequences_key" ON "public"."sequences" ("key");
-- Create index "not null" to table: "sequences"
CREATE UNIQUE INDEX "not null" ON "public"."sequences" ("idempotency_key");
-- Create "doi_deposits" table
CREATE TABLE "public"."doi_deposits" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "article_id" bigint NOT NULL,
  "batch_id" text NOT NULL,
  "doi" text NOT NULL,
  "status" text NOT NULL DEFAULT 'PENDING',
  "message" text NULL,
  "submission_id" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_doi_deposits_article" FOREIGN KEY ("article_id") REFERENCES "public"."articles" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_doi_deposits_article_id" to table: "doi_deposits"
CREATE INDEX "idx_doi_deposits_article_id" ON "public"."doi_deposits" ("article_id");
-- Create index "idx_doi_deposits_batch_id" to table: "doi_deposits"
CREATE UNIQUE INDEX "idx_doi_deposits_batch_id" ON "public"."doi_deposits" ("batch_id");
-- Create index "idx_doi_deposits_deleted_at" to table: "doi_deposits"
CREATE INDEX "idx_doi_deposits_deleted_at" ON "public"."doi_deposits" ("deleted_at");
