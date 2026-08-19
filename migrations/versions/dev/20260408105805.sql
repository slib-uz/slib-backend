-- Create "instructions" table
CREATE TABLE "public"."instructions" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "key" character varying(64) NOT NULL,
  "video_link" character varying(1024) NULL,
  "description" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_instructions_deleted_at" to table: "instructions"
CREATE INDEX "idx_instructions_deleted_at" ON "public"."instructions" ("deleted_at");
-- Create index "idx_instructions_key" to table: "instructions"
CREATE UNIQUE INDEX "idx_instructions_key" ON "public"."instructions" ("key");
-- Create "journal_news" table
CREATE TABLE "public"."journal_news" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "journal_id" bigint NOT NULL,
  "title" jsonb NULL,
  "body" jsonb NULL,
  "image" character varying(1024) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_journal_news_journal" FOREIGN KEY ("journal_id") REFERENCES "public"."journals" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_journal_news_deleted_at" to table: "journal_news"
CREATE INDEX "idx_journal_news_deleted_at" ON "public"."journal_news" ("deleted_at");
-- Create index "idx_journal_news_journal_id" to table: "journal_news"
CREATE INDEX "idx_journal_news_journal_id" ON "public"."journal_news" ("journal_id");
