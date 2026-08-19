-- Modify "roles" table
ALTER TABLE "public"."roles" ADD COLUMN "institution_id" bigint NULL, ADD CONSTRAINT "fk_roles_institution" FOREIGN KEY ("institution_id") REFERENCES "public"."institutions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Create index "idx_roles_institution_id" to table: "roles"
CREATE INDEX "idx_roles_institution_id" ON "public"."roles" ("institution_id");
