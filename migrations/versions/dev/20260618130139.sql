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
ALTER TABLE "public"."publishers" ADD COLUMN "institution_id" bigint NULL, ADD CONSTRAINT "fk_publishers_institution" FOREIGN KEY ("institution_id") REFERENCES "public"."institutions" ("id") ON UPDATE NO ACTION ON DELETE SET NULL;
-- Create index "idx_publishers_institution_id" to table: "publishers"
CREATE INDEX "idx_publishers_institution_id" ON "public"."publishers" ("institution_id");
