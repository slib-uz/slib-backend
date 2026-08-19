-- Create "legacy_authors" table
CREATE TABLE "public"."legacy_authors" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "full_name" character varying(512) NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_legacy_authors_deleted_at" to table: "legacy_authors"
CREATE INDEX "idx_legacy_authors_deleted_at" ON "public"."legacy_authors" ("deleted_at");
-- Create index "idx_legacy_authors_full_name" to table: "legacy_authors"
CREATE INDEX "idx_legacy_authors_full_name" ON "public"."legacy_authors" ("full_name");
-- Create "legacy_author_articles" table
CREATE TABLE "public"."legacy_author_articles" (
  "legacy_author_model_id" bigint NOT NULL,
  "article_model_id" bigint NOT NULL,
  PRIMARY KEY ("legacy_author_model_id", "article_model_id"),
  CONSTRAINT "fk_legacy_author_articles_article_model" FOREIGN KEY ("article_model_id") REFERENCES "public"."articles" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_legacy_author_articles_legacy_author_model" FOREIGN KEY ("legacy_author_model_id") REFERENCES "public"."legacy_authors" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
