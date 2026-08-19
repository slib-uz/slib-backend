-- Create "instructions" table
CREATE TABLE "public"."instructions" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "key" character varying(64) NOT NULL,
  "video_link" character varying(1024) NULL,
  "description" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_instructions_deleted_at" to table: "instructions"
CREATE INDEX "idx_instructions_deleted_at" ON "public"."instructions" ("deleted_at");
-- Create index "idx_key_deleted_at" to table: "instructions"
CREATE UNIQUE INDEX "idx_key_deleted_at" ON "public"."instructions" ("key") WHERE (deleted_at IS NULL);
-- Create "doi_deposits" table
CREATE TABLE "public"."doi_deposits" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "article_id" bigint NULL,
  "batch_id" text NOT NULL,
  "doi" text NOT NULL,
  "status" text NOT NULL DEFAULT 'PENDING',
  "message" text NULL,
  "submission_id" text NULL,
  "request_body" text NULL,
  "response_body" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_doi_deposits_article" FOREIGN KEY ("article_id") REFERENCES "public"."articles" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create index "idx_doi_deposits_article_id" to table: "doi_deposits"
CREATE INDEX "idx_doi_deposits_article_id" ON "public"."doi_deposits" ("article_id");
-- Create index "idx_doi_deposits_batch_id" to table: "doi_deposits"
CREATE UNIQUE INDEX "idx_doi_deposits_batch_id" ON "public"."doi_deposits" ("batch_id");
-- Create index "idx_doi_deposits_deleted_at" to table: "doi_deposits"
CREATE INDEX "idx_doi_deposits_deleted_at" ON "public"."doi_deposits" ("deleted_at");
-- Modify "article_purchases" table
ALTER TABLE "public"."article_purchases" DROP COLUMN "invoice_id";
-- Modify "journal_configs" table
ALTER TABLE "public"."journal_configs" ADD COLUMN "creator_id" bigint NULL, ADD CONSTRAINT "fk_journal_configs_creator" FOREIGN KEY ("creator_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL;
-- Modify "journals" table
ALTER TABLE "public"."journals" DROP COLUMN "language_id", ADD COLUMN "submission_access" bigint NOT NULL DEFAULT 10, ADD COLUMN "comment_access" bigint NOT NULL DEFAULT 10, ADD COLUMN "social_networks" jsonb NULL, ADD COLUMN "support_link" character varying(255) NULL;
-- Create "journal_doi_settings" table
CREATE TABLE "public"."journal_doi_settings" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "journal_id" bigint NOT NULL,
  "username" text NOT NULL,
  "password" text NOT NULL,
  "doi_prefix" text NOT NULL,
  "doi_suffix" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_journal_doi_settings_journal" FOREIGN KEY ("journal_id") REFERENCES "public"."journals" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_journal_doi_settings_deleted_at" to table: "journal_doi_settings"
CREATE INDEX "idx_journal_doi_settings_deleted_at" ON "public"."journal_doi_settings" ("deleted_at");
-- Create index "idx_journal_doi_settings_journal_id" to table: "journal_doi_settings"
CREATE UNIQUE INDEX "idx_journal_doi_settings_journal_id" ON "public"."journal_doi_settings" ("journal_id");
-- Create "journal_languages" table
CREATE TABLE "public"."journal_languages" (
  "journal_model_id" bigint NOT NULL,
  "language_model_id" bigint NOT NULL,
  PRIMARY KEY ("journal_model_id", "language_model_id"),
  CONSTRAINT "fk_journal_languages_journal_model" FOREIGN KEY ("journal_model_id") REFERENCES "public"."journals" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_journal_languages_language_model" FOREIGN KEY ("language_model_id") REFERENCES "public"."languages" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "journal_news" table
CREATE TABLE "public"."journal_news" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "journal_id" bigint NOT NULL,
  "title" jsonb NULL,
  "body" jsonb NULL,
  "image" character varying(1024) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_journal_news_journal" FOREIGN KEY ("journal_id") REFERENCES "public"."journals" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_journal_news_deleted_at" to table: "journal_news"
CREATE INDEX "idx_journal_news_deleted_at" ON "public"."journal_news" ("deleted_at");
-- Create index "idx_journal_news_journal_id" to table: "journal_news"
CREATE INDEX "idx_journal_news_journal_id" ON "public"."journal_news" ("journal_id");
-- Drop "article_requirements" table
DROP TABLE "public"."article_requirements";
-- Drop "transactions" table
DROP TABLE "public"."transactions";
-- Drop "invoices" table
DROP TABLE "public"."invoices";
