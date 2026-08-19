-- Create "journal_editorials" table
CREATE TABLE "public"."journal_editorials" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "journal_id" bigint NOT NULL,
  "full_name" character varying(128) NOT NULL,
  "role_title" character varying(128) NOT NULL,
  "photo" character varying(128) NULL,
  "science_id" character varying(128) NULL,
  "workplace" character varying(128) NULL,
  "position" character varying(128) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_journal_editorials_journal" FOREIGN KEY ("journal_id") REFERENCES "public"."journals" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_journal_editorials_deleted_at" to table: "journal_editorials"
CREATE INDEX "idx_journal_editorials_deleted_at" ON "public"."journal_editorials" ("deleted_at");
-- Create index "idx_journal_editorials_journal_id" to table: "journal_editorials"
CREATE INDEX "idx_journal_editorials_journal_id" ON "public"."journal_editorials" ("journal_id");
