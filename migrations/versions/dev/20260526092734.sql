-- Modify "journal_doi_settings" table
ALTER TABLE "public"."journal_doi_settings" ADD COLUMN "journal_name" text NULL;
-- Modify "ai_detect_result_models" table
ALTER TABLE "public"."ai_detect_result_models" DROP CONSTRAINT "fk_ai_detect_result_models_application", DROP CONSTRAINT "fk_ai_detect_result_models_review_stage", ADD CONSTRAINT "fk_article_applications_ai_detect_results" FOREIGN KEY ("application_id") REFERENCES "public"."article_applications" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "fk_review_stages_ai_detect_results" FOREIGN KEY ("review_stage_id") REFERENCES "public"."review_stages" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
