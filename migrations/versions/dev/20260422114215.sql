-- Create "journal_doi_settings" table
CREATE TABLE "public"."journal_doi_settings" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "journal_id" bigint NOT NULL,
  "username" text NOT NULL,
  "password" text NOT NULL,
  "do_iprefix" text NOT NULL,
  "doi_suffix" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_journal_doi_settings_journal" FOREIGN KEY ("journal_id") REFERENCES "public"."journals" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_journal_doi_settings_deleted_at" to table: "journal_doi_settings"
CREATE INDEX "idx_journal_doi_settings_deleted_at" ON "public"."journal_doi_settings" ("deleted_at");
-- Create index "idx_journal_doi_settings_journal_id" to table: "journal_doi_settings"
CREATE UNIQUE INDEX "idx_journal_doi_settings_journal_id" ON "public"."journal_doi_settings" ("journal_id");
