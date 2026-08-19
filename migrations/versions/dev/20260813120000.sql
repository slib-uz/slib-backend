-- Modify "tags" table
ALTER TABLE "public"."tags"
  ADD COLUMN "created_at" timestamptz NULL,
  ADD COLUMN "updated_at" timestamptz NULL,
  ADD COLUMN "deleted_at" timestamptz NULL,
  ADD COLUMN "slug" character varying(64) NULL;
-- Create "tag_translations" table
CREATE TABLE "public"."tag_translations" (
  "id" bigserial NOT NULL,
  "tag_id" bigint NOT NULL,
  "language_code" character varying(5) NOT NULL,
  "name" character varying(32) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_tag_translations_tag" FOREIGN KEY ("tag_id") REFERENCES "public"."tags" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Migrate existing tag names as Uzbek translations
INSERT INTO "public"."tag_translations" ("tag_id", "language_code", "name")
SELECT "id", 'uz', "name" FROM "public"."tags" WHERE "name" IS NOT NULL AND "name" <> '';
-- Backfill slugs for existing tags
UPDATE "public"."tags" SET "slug" = 'tag-' || "id" WHERE "slug" IS NULL;
-- Set slug required and drop old name column
ALTER TABLE "public"."tags" ALTER COLUMN "slug" SET NOT NULL, DROP COLUMN "name";
-- Create index "idx_tags_deleted_at" to table: "tags"
CREATE INDEX "idx_tags_deleted_at" ON "public"."tags" ("deleted_at");
-- Create index "idx_tags_slug" to table: "tags"
CREATE UNIQUE INDEX "idx_tags_slug" ON "public"."tags" ("slug") WHERE (deleted_at IS NULL);
-- Create index "idx_tag_tr_tag_lang" to table: "tag_translations"
CREATE UNIQUE INDEX "idx_tag_tr_tag_lang" ON "public"."tag_translations" ("tag_id", "language_code");
-- Create index "idx_tag_translations_name" to table: "tag_translations"
CREATE INDEX "idx_tag_translations_name" ON "public"."tag_translations" ("name");
-- Create index "idx_tag_translations_lang_name" to table: "tag_translations"
CREATE UNIQUE INDEX "idx_tag_translations_lang_name" ON "public"."tag_translations" ("language_code", "name");
