-- Drop index "idx_instructions_key" from table: "instructions"
DROP INDEX "public"."idx_instructions_key";
-- Create index "idx_key_deleted_at" to table: "instructions"
CREATE UNIQUE INDEX "idx_key_deleted_at" ON "public"."instructions" ("key") WHERE (deleted_at IS NULL);
