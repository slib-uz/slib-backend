-- Modify "user_profiles" table
ALTER TABLE "public"."user_profiles" ADD COLUMN "last_online_at" timestamptz NULL;
-- Create "regions" table
CREATE TABLE "public"."regions" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" jsonb NOT NULL,
  "soato" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_regions_soato" UNIQUE ("soato")
);
-- Create index "idx_regions_deleted_at" to table: "regions"
CREATE INDEX "idx_regions_deleted_at" ON "public"."regions" ("deleted_at");
-- Create "sandbox_users" table
CREATE TABLE "public"."sandbox_users" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "full_name" text NULL,
  "science_id" text NULL,
  "phone_number" text NULL,
  "otp" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_sandbox_users_deleted_at" to table: "sandbox_users"
CREATE INDEX "idx_sandbox_users_deleted_at" ON "public"."sandbox_users" ("deleted_at");
-- Create "districts" table
CREATE TABLE "public"."districts" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" jsonb NOT NULL,
  "soato" bigint NOT NULL,
  "region_id" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_districts_soato" UNIQUE ("soato"),
  CONSTRAINT "fk_districts_region" FOREIGN KEY ("region_id") REFERENCES "public"."regions" ("id") ON UPDATE CASCADE ON DELETE SET NULL
);
-- Create index "idx_districts_deleted_at" to table: "districts"
CREATE INDEX "idx_districts_deleted_at" ON "public"."districts" ("deleted_at");
-- Modify "journals" table
ALTER TABLE "public"."journals" ADD COLUMN "region_id" bigint NULL, ADD COLUMN "district_id" bigint NULL, ADD CONSTRAINT "fk_journals_district" FOREIGN KEY ("district_id") REFERENCES "public"."districts" ("id") ON UPDATE NO ACTION ON DELETE SET NULL, ADD CONSTRAINT "fk_journals_region" FOREIGN KEY ("region_id") REFERENCES "public"."regions" ("id") ON UPDATE NO ACTION ON DELETE SET NULL;
-- Create index "idx_journals_district_id" to table: "journals"
CREATE INDEX "idx_journals_district_id" ON "public"."journals" ("district_id");
-- Create index "idx_journals_region_id" to table: "journals"
CREATE INDEX "idx_journals_region_id" ON "public"."journals" ("region_id");
-- Create "institutions" table
CREATE TABLE "public"."institutions" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" character varying(500) NOT NULL,
  "tin" character varying(32) NULL,
  "logo" character varying(500) NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_institutions_deleted_at" to table: "institutions"
CREATE INDEX "idx_institutions_deleted_at" ON "public"."institutions" ("deleted_at");
-- Modify "publishers" table
ALTER TABLE "public"."publishers" ADD COLUMN "institution_id" bigint NULL, ADD COLUMN "region_id" bigint NULL, ADD COLUMN "district_id" bigint NULL, ADD CONSTRAINT "fk_institutions_publishers" FOREIGN KEY ("institution_id") REFERENCES "public"."institutions" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_publishers_district" FOREIGN KEY ("district_id") REFERENCES "public"."districts" ("id") ON UPDATE NO ACTION ON DELETE SET NULL, ADD CONSTRAINT "fk_publishers_region" FOREIGN KEY ("region_id") REFERENCES "public"."regions" ("id") ON UPDATE NO ACTION ON DELETE SET NULL;
-- Create index "idx_publishers_district_id" to table: "publishers"
CREATE INDEX "idx_publishers_district_id" ON "public"."publishers" ("district_id");
-- Create index "idx_publishers_institution_id" to table: "publishers"
CREATE INDEX "idx_publishers_institution_id" ON "public"."publishers" ("institution_id");
-- Create index "idx_publishers_region_id" to table: "publishers"
CREATE INDEX "idx_publishers_region_id" ON "public"."publishers" ("region_id");
-- Modify "roles" table
ALTER TABLE "public"."roles" ADD COLUMN "institution_id" bigint NULL, ADD CONSTRAINT "fk_roles_institution" FOREIGN KEY ("institution_id") REFERENCES "public"."institutions" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-- Create index "idx_roles_institution_id" to table: "roles"
CREATE INDEX "idx_roles_institution_id" ON "public"."roles" ("institution_id");
