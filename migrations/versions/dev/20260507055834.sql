-- Create "outbox_events" table
CREATE TABLE "public"."outbox_events" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "event_id" character varying(50) NOT NULL,
  "version" bigint NOT NULL,
  "event" character varying(50) NOT NULL,
  "payload" jsonb NULL,
  "delivered_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_outbox_events_deleted_at" to table: "outbox_events"
CREATE INDEX "idx_outbox_events_deleted_at" ON "public"."outbox_events" ("deleted_at");
-- Create index "idx_outbox_events_event_id" to table: "outbox_events"
CREATE UNIQUE INDEX "idx_outbox_events_event_id" ON "public"."outbox_events" ("event_id");
-- Create "ai_detect_result_models" table
CREATE TABLE "public"."ai_detect_result_models" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "review_stage_id" bigint NULL,
  "application_id" bigint NULL,
  "article_id" bigint NULL,
  "journal_id" bigint NULL,
  "external_id" bigint NOT NULL,
  "words_count" bigint NOT NULL DEFAULT 0,
  "status" bigint NOT NULL,
  "status_display" character varying(128) NULL,
  "human_percent" numeric NOT NULL DEFAULT 0,
  "warn_percent" numeric NOT NULL DEFAULT 0,
  "ai_percent" numeric NOT NULL DEFAULT 0,
  "report_url" character varying(2048) NULL,
  "external_created_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_ai_detect_result_models_application" FOREIGN KEY ("application_id") REFERENCES "public"."article_applications" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "fk_ai_detect_result_models_article" FOREIGN KEY ("article_id") REFERENCES "public"."articles" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "fk_ai_detect_result_models_journal" FOREIGN KEY ("journal_id") REFERENCES "public"."journals" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "fk_ai_detect_result_models_review_stage" FOREIGN KEY ("review_stage_id") REFERENCES "public"."review_stages" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create index "idx_ai_detect_result_models_deleted_at" to table: "ai_detect_result_models"
CREATE INDEX "idx_ai_detect_result_models_deleted_at" ON "public"."ai_detect_result_models" ("deleted_at");
-- Create index "idx_ai_detect_result_models_external_id" to table: "ai_detect_result_models"
CREATE UNIQUE INDEX "idx_ai_detect_result_models_external_id" ON "public"."ai_detect_result_models" ("external_id");
