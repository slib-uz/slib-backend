-- Drop index "idx_article_id_author_id" from table: "article_author_affiliations"
DROP INDEX "public"."idx_article_id_author_id";
-- Create index "idx_organizations_tin" to table: "organizations"
CREATE UNIQUE INDEX "idx_organizations_tin" ON "public"."organizations" ("tin") WHERE (deleted_at IS NULL);
-- Drop index "idx_tags_name" from table: "tags"
DROP INDEX "public"."idx_tags_name";
-- Modify "tags" table
ALTER TABLE "public"."tags" ADD COLUMN "lang" character varying(5) NOT NULL DEFAULT 'uz';
-- Create index "idx_tags_lang_name" to table: "tags"
CREATE UNIQUE INDEX "idx_tags_lang_name" ON "public"."tags" ("name", "lang");
