-- Modify "journals" table
ALTER TABLE "public"."journals" DROP COLUMN "language_id", ADD COLUMN "submission_access" bigint NOT NULL DEFAULT 10, ADD COLUMN "comment_access" bigint NOT NULL DEFAULT 10, ADD COLUMN "social_networks" jsonb NULL, ADD COLUMN "support_link" character varying(255) NULL;
-- Create "journal_languages" table
CREATE TABLE "public"."journal_languages" (
  "journal_model_id" bigint NOT NULL,
  "language_model_id" bigint NOT NULL,
  PRIMARY KEY ("journal_model_id", "language_model_id"),
  CONSTRAINT "fk_journal_languages_journal_model" FOREIGN KEY ("journal_model_id") REFERENCES "public"."journals" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_journal_languages_language_model" FOREIGN KEY ("language_model_id") REFERENCES "public"."languages" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Drop "article_requirements" table
DROP TABLE "public"."article_requirements";
