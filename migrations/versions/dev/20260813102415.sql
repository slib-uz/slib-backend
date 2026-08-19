-- Normalize STIR/TIN to digits only before unique index
UPDATE "public"."organizations"
SET "tin" = regexp_replace("tin", '[^0-9]', '', 'g')
WHERE "tin" ~ '[^0-9]';

-- Create index "idx_organizations_tin" to table: "organizations"
CREATE UNIQUE INDEX "idx_organizations_tin" ON "public"."organizations" ("tin") WHERE (deleted_at IS NULL);
