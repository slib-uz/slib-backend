-- Modify "journals" table
ALTER TABLE "public"."journals" ADD COLUMN "region_id" bigint NULL, ADD COLUMN "district_id" bigint NULL, ADD CONSTRAINT "fk_journals_district" FOREIGN KEY ("district_id") REFERENCES "public"."districts" ("id") ON UPDATE NO ACTION ON DELETE SET NULL, ADD CONSTRAINT "fk_journals_region" FOREIGN KEY ("region_id") REFERENCES "public"."regions" ("id") ON UPDATE NO ACTION ON DELETE SET NULL;
-- Create index "idx_journals_district_id" to table: "journals"
CREATE INDEX "idx_journals_district_id" ON "public"."journals" ("district_id");
-- Create index "idx_journals_region_id" to table: "journals"
CREATE INDEX "idx_journals_region_id" ON "public"."journals" ("region_id");
-- Modify "publishers" table
ALTER TABLE "public"."publishers" DROP CONSTRAINT "fk_publishers_institution", ADD COLUMN "region_id" bigint NULL, ADD COLUMN "district_id" bigint NULL, ADD CONSTRAINT "fk_institutions_publishers" FOREIGN KEY ("institution_id") REFERENCES "public"."institutions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_publishers_district" FOREIGN KEY ("district_id") REFERENCES "public"."districts" ("id") ON UPDATE NO ACTION ON DELETE SET NULL, ADD CONSTRAINT "fk_publishers_region" FOREIGN KEY ("region_id") REFERENCES "public"."regions" ("id") ON UPDATE NO ACTION ON DELETE SET NULL;
-- Create index "idx_publishers_district_id" to table: "publishers"
CREATE INDEX "idx_publishers_district_id" ON "public"."publishers" ("district_id");
-- Create index "idx_publishers_region_id" to table: "publishers"
CREATE INDEX "idx_publishers_region_id" ON "public"."publishers" ("region_id");
