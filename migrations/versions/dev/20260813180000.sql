-- Restore tags.name and add lang; split translations into independent tag rows
ALTER TABLE "public"."tags"
  ADD COLUMN "name" character varying(32) NULL,
  ADD COLUMN "lang" character varying(5) NOT NULL DEFAULT 'uz',
  ALTER COLUMN "slug" DROP NOT NULL;
-- Copy preferred translation (uz, then en, then other) onto existing tag rows
UPDATE "public"."tags" AS t
SET
  "name" = sub."name",
  "lang" = sub."language_code"
FROM (
  SELECT DISTINCT ON ("tag_id") "tag_id", "name", "language_code"
  FROM "public"."tag_translations"
  ORDER BY "tag_id",
    CASE "language_code"
      WHEN 'uz' THEN 0
      WHEN 'en' THEN 1
      ELSE 2
    END
) AS sub
WHERE t."id" = sub."tag_id";
-- Fallback for tags with no translations
UPDATE "public"."tags"
SET "name" = LEFT("slug", 32)
WHERE "name" IS NULL OR "name" = '';
-- Remaining translations become new independent tags
INSERT INTO "public"."tags" ("name", "lang")
SELECT DISTINCT tt."name", tt."language_code"
FROM "public"."tag_translations" tt
INNER JOIN "public"."tags" t ON t."id" = tt."tag_id"
WHERE NOT (tt."language_code" = t."lang" AND tt."name" = t."name");
-- Attach the new per-language tags to the same articles
INSERT INTO "public"."article_tags" ("article_model_id", "tag_model_id")
SELECT DISTINCT at."article_model_id", nt."id"
FROM "public"."article_tags" at
INNER JOIN "public"."tag_translations" tt ON tt."tag_id" = at."tag_model_id"
INNER JOIN "public"."tags" nt ON nt."name" = tt."name" AND nt."lang" = tt."language_code"
WHERE nt."id" <> at."tag_model_id"
ON CONFLICT ("article_model_id", "tag_model_id") DO NOTHING;
-- Require name and unique (lang, name)
ALTER TABLE "public"."tags" ALTER COLUMN "name" SET NOT NULL;
CREATE UNIQUE INDEX "idx_tags_lang_name" ON "public"."tags" ("lang", "name");
-- Drop translation table and slug/timestamps
DROP TABLE "public"."tag_translations";
DROP INDEX IF EXISTS "public"."idx_tags_slug";
DROP INDEX IF EXISTS "public"."idx_tags_deleted_at";
ALTER TABLE "public"."tags"
  DROP COLUMN "slug",
  DROP COLUMN "created_at",
  DROP COLUMN "updated_at",
  DROP COLUMN "deleted_at";
