-- Drop index "idx_tags_lang_name" from table: "tags"
DROP INDEX "public"."idx_tags_lang_name";
-- Create index "idx_tags_lang_name" to table: "tags"
CREATE UNIQUE INDEX "idx_tags_lang_name" ON "public"."tags" ("name", "lang");
