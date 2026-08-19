-- Create "notification_send_result_models" table
CREATE TABLE "public"."notification_send_result_models" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "notification_id" bigint NOT NULL,
  "success_count" bigint NOT NULL DEFAULT 0,
  "failure_count" bigint NOT NULL DEFAULT 0,
  "failed_tokens" jsonb NULL,
  "errors" jsonb NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_notification_send_result_models_deleted_at" to table: "notification_send_result_models"
CREATE INDEX "idx_notification_send_result_models_deleted_at" ON "public"."notification_send_result_models" ("deleted_at");
-- Create index "idx_notification_send_result_models_notification_id" to table: "notification_send_result_models"
CREATE INDEX "idx_notification_send_result_models_notification_id" ON "public"."notification_send_result_models" ("notification_id");
-- Create "task_result_models" table
CREATE TABLE "public"."task_result_models" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "task_id" text NOT NULL,
  "task" character varying(32) NOT NULL,
  "payload" text NULL,
  "status" character varying(32) NULL DEFAULT 'PENDING',
  "error" text NULL,
  "retry_count" bigint NULL,
  "completed_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_task_result_models_deleted_at" to table: "task_result_models"
CREATE INDEX "idx_task_result_models_deleted_at" ON "public"."task_result_models" ("deleted_at");
-- Create index "idx_task_result_models_task_id" to table: "task_result_models"
CREATE INDEX "idx_task_result_models_task_id" ON "public"."task_result_models" ("task_id");
-- Create "ethics_policies" table
CREATE TABLE "public"."ethics_policies" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "content" jsonb NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_ethics_policies_deleted_at" to table: "ethics_policies"
CREATE INDEX "idx_ethics_policies_deleted_at" ON "public"."ethics_policies" ("deleted_at");
-- Create "public_offers" table
CREATE TABLE "public"."public_offers" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "description" jsonb NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_public_offers_deleted_at" to table: "public_offers"
CREATE INDEX "idx_public_offers_deleted_at" ON "public"."public_offers" ("deleted_at");
-- Create "projects" table
CREATE TABLE "public"."projects" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "title" character varying(512) NULL,
  "logo_path" character varying(512) NULL,
  "link" character varying(512) NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_projects_deleted_at" to table: "projects"
CREATE INDEX "idx_projects_deleted_at" ON "public"."projects" ("deleted_at");
-- Create "partners" table
CREATE TABLE "public"."partners" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "title" character varying(512) NULL,
  "logo_path" character varying(512) NULL,
  "link" character varying(512) NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_partners_deleted_at" to table: "partners"
CREATE INDEX "idx_partners_deleted_at" ON "public"."partners" ("deleted_at");
-- Create "otp_codes" table
CREATE TABLE "public"."otp_codes" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "phone" character varying(50) NOT NULL,
  "code" character varying(32) NOT NULL,
  "session_id" character varying(255) NOT NULL,
  "purpose" character varying(32) NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "used_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_otp_codes_deleted_at" to table: "otp_codes"
CREATE INDEX "idx_otp_codes_deleted_at" ON "public"."otp_codes" ("deleted_at");
-- Create "notification_template_models" table
CREATE TABLE "public"."notification_template_models" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "key" character varying(64) NOT NULL,
  "title" jsonb NULL,
  "body" jsonb NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_notification_template_models_deleted_at" to table: "notification_template_models"
CREATE INDEX "idx_notification_template_models_deleted_at" ON "public"."notification_template_models" ("deleted_at");
-- Create index "idx_notification_template_models_key" to table: "notification_template_models"
CREATE UNIQUE INDEX "idx_notification_template_models_key" ON "public"."notification_template_models" ("key");
-- Create "guides" table
CREATE TABLE "public"."guides" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "title" jsonb NOT NULL,
  "description" jsonb NOT NULL,
  "file_path" text NULL,
  "video_url" text NULL DEFAULT '',
  PRIMARY KEY ("id")
);
-- Create index "idx_guides_deleted_at" to table: "guides"
CREATE INDEX "idx_guides_deleted_at" ON "public"."guides" ("deleted_at");
-- Create "abouts" table
CREATE TABLE "public"."abouts" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "content" jsonb NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_abouts_deleted_at" to table: "abouts"
CREATE INDEX "idx_abouts_deleted_at" ON "public"."abouts" ("deleted_at");
-- Create "users" table
CREATE TABLE "public"."users" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "science_id" character varying(64) NOT NULL,
  "pin" character varying(64) NULL,
  "first_name" character varying(255) NULL,
  "last_name" character varying(255) NULL,
  "middle_name" character varying(255) NULL,
  "full_name" character varying(255) NULL,
  "gender" bigint NULL,
  "birth_date" date NULL,
  "phone_number" character varying(64) NOT NULL,
  "is_admin" boolean NULL DEFAULT false,
  "photo" character varying(255) NULL,
  "email" character varying(255) NULL,
  "city" character varying(5) NOT NULL DEFAULT '',
  "academic_degree" character varying(5) NOT NULL DEFAULT '',
  "academic_title" character varying(255) NULL,
  "orc_id_id" character varying(50) NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "public"."users" ("deleted_at");
-- Create index "idx_users_orc_id_id" to table: "users"
CREATE UNIQUE INDEX "idx_users_orc_id_id" ON "public"."users" ("orc_id_id");
-- Create index "idx_users_phone_number" to table: "users"
CREATE UNIQUE INDEX "idx_users_phone_number" ON "public"."users" ("phone_number");
-- Create index "idx_users_pin" to table: "users"
CREATE UNIQUE INDEX "idx_users_pin" ON "public"."users" ("pin");
-- Create index "idx_users_science_id" to table: "users"
CREATE UNIQUE INDEX "idx_users_science_id" ON "public"."users" ("science_id");
-- Create "academic_degrees" table
CREATE TABLE "public"."academic_degrees" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "source_id" bigint NOT NULL,
  "speciality" text NULL,
  "confirmed_date" date NULL,
  "user_id" bigint NOT NULL,
  "diploma_number" text NULL,
  "science_sector" text NULL,
  "degree_name" text NULL,
  "degree_code" bigint NULL,
  "science_sector_code" bigint NULL,
  "theme" text NULL,
  "awarded_at" date NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_academic_degrees_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_academic_degrees_deleted_at" to table: "academic_degrees"
CREATE INDEX "idx_academic_degrees_deleted_at" ON "public"."academic_degrees" ("deleted_at");
-- Create index "idx_academic_degrees_source_id" to table: "academic_degrees"
CREATE INDEX "idx_academic_degrees_source_id" ON "public"."academic_degrees" ("source_id");
-- Create "academic_titles" table
CREATE TABLE "public"."academic_titles" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "source_id" bigint NOT NULL,
  "title" text NULL,
  "confirmed_date" date NULL,
  "diploma_number" text NULL,
  "user_id" bigint NOT NULL,
  "science_sector" text NULL,
  "science_sector_code" bigint NULL,
  "title_code" bigint NULL,
  "speciality" text NULL,
  "awarded_at" date NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_academic_titles_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_academic_titles_deleted_at" to table: "academic_titles"
CREATE INDEX "idx_academic_titles_deleted_at" ON "public"."academic_titles" ("deleted_at");
-- Create index "idx_academic_titles_source_id" to table: "academic_titles"
CREATE INDEX "idx_academic_titles_source_id" ON "public"."academic_titles" ("source_id");
-- Create "languages" table
CREATE TABLE "public"."languages" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" jsonb NULL,
  "code" character varying(32) NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_languages_deleted_at" to table: "languages"
CREATE INDEX "idx_languages_deleted_at" ON "public"."languages" ("deleted_at");
-- Create "publishers" table
CREATE TABLE "public"."publishers" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "tin" character varying(32) NULL,
  "name" character varying(500) NULL,
  "short_name" character varying(255) NULL,
  "logo" character varying(500) NULL,
  "phone_number" character varying(128) NULL,
  "fax_number" character varying(128) NULL,
  "email" character varying(255) NULL,
  "website" character varying(64) NULL,
  "address" character varying(500) NULL,
  "description" character varying(32512) NULL,
  "is_active" boolean NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_publishers_deleted_at" to table: "publishers"
CREATE INDEX "idx_publishers_deleted_at" ON "public"."publishers" ("deleted_at");
-- Create index "idx_publishers_tin" to table: "publishers"
CREATE UNIQUE INDEX "idx_publishers_tin" ON "public"."publishers" ("tin");
-- Create "journals" table
CREATE TABLE "public"."journals" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" jsonb NULL,
  "short_name" jsonb NULL,
  "issn_paper" character varying(255) NULL,
  "issn_online" character varying(255) NULL,
  "description" jsonb NULL,
  "rule" jsonb NULL,
  "article_publish_conditions" jsonb NULL,
  "established_date" timestamptz NULL,
  "website" character varying(128) NULL,
  "certificate_file" character varying(500) NULL,
  "email" character varying(255) NULL,
  "address" jsonb NULL,
  "phone_number" character varying(255) NULL,
  "cover_image_file" character varying(500) NULL,
  "publisher_id" bigint NOT NULL,
  "publishing_price" bigint NULL,
  "selling_price" bigint NULL,
  "access_type" bigint NULL,
  "oak_certificate_file" character varying(500) NULL,
  "language_id" bigint NULL,
  "peer_review_method" bigint NOT NULL,
  "is_active" boolean NOT NULL DEFAULT false,
  "is_accepted" boolean NOT NULL DEFAULT false,
  "rating_count" bigint NOT NULL DEFAULT 0,
  "rating_sum" bigint NOT NULL DEFAULT 0,
  "views_count" bigint NOT NULL DEFAULT 0,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_journals_language" FOREIGN KEY ("language_id") REFERENCES "public"."languages" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "fk_journals_publisher" FOREIGN KEY ("publisher_id") REFERENCES "public"."publishers" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create index "idx_journals_deleted_at" to table: "journals"
CREATE INDEX "idx_journals_deleted_at" ON "public"."journals" ("deleted_at");
-- Create index "idx_journals_issn_online" to table: "journals"
CREATE INDEX "idx_journals_issn_online" ON "public"."journals" ("issn_online");
-- Create index "idx_journals_issn_paper" to table: "journals"
CREATE INDEX "idx_journals_issn_paper" ON "public"."journals" ("issn_paper");
-- Create index "idx_journals_publisher_id" to table: "journals"
CREATE INDEX "idx_journals_publisher_id" ON "public"."journals" ("publisher_id");
-- Create "editions" table
CREATE TABLE "public"."editions" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "journal_id" bigint NULL,
  "file" character varying(1024) NULL,
  "cover_image" character varying(1024) NULL,
  "name" character varying(255) NULL,
  "number" character varying(255) NULL,
  "volume" character varying(255) NULL,
  "published_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_editions_journal" FOREIGN KEY ("journal_id") REFERENCES "public"."journals" ("id") ON UPDATE CASCADE ON DELETE SET NULL
);
-- Create index "idx_editions_deleted_at" to table: "editions"
CREATE INDEX "idx_editions_deleted_at" ON "public"."editions" ("deleted_at");
-- Create "articles" table
CREATE TABLE "public"."articles" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" jsonb NULL,
  "publication_date" date NULL,
  "co_authors_count" bigint NOT NULL,
  "access_type" bigint NOT NULL,
  "language_id" bigint NOT NULL,
  "annotation" jsonb NULL,
  "content_file_path" text NULL,
  "doi" character varying(255) NULL,
  "roi" character varying(255) NULL,
  "journal_id" bigint NULL,
  "expert_conclusion_file" text NULL,
  "is_published" boolean NOT NULL DEFAULT false,
  "views_count" bigint NOT NULL DEFAULT 0,
  "rating_count" bigint NOT NULL DEFAULT 0,
  "rating_sum" bigint NOT NULL DEFAULT 0,
  "external_id" bigint NULL,
  "is_sent" boolean NOT NULL DEFAULT false,
  "edition_id" bigint NULL,
  "pages" character varying(255) NULL,
  "web_address" character varying(512) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_articles_edition" FOREIGN KEY ("edition_id") REFERENCES "public"."editions" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "fk_articles_journal" FOREIGN KEY ("journal_id") REFERENCES "public"."journals" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "fk_articles_language" FOREIGN KEY ("language_id") REFERENCES "public"."languages" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create index "idx_articles_deleted_at" to table: "articles"
CREATE INDEX "idx_articles_deleted_at" ON "public"."articles" ("deleted_at");
-- Create index "idx_articles_doi" to table: "articles"
CREATE UNIQUE INDEX "idx_articles_doi" ON "public"."articles" ("doi");
-- Create index "idx_articles_external_id" to table: "articles"
CREATE UNIQUE INDEX "idx_articles_external_id" ON "public"."articles" ("external_id");
-- Create index "idx_articles_roi" to table: "articles"
CREATE UNIQUE INDEX "idx_articles_roi" ON "public"."articles" ("roi");
-- Create "article_applications" table
CREATE TABLE "public"."article_applications" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "number" character varying(255) NULL,
  "article_id" bigint NULL,
  "journal_id" bigint NULL,
  "user_id" bigint NULL,
  "is_published" boolean NOT NULL DEFAULT false,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_article_applications_article" FOREIGN KEY ("article_id") REFERENCES "public"."articles" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_article_applications_journal" FOREIGN KEY ("journal_id") REFERENCES "public"."journals" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "fk_article_applications_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create index "idx_article_applications_article_id" to table: "article_applications"
CREATE INDEX "idx_article_applications_article_id" ON "public"."article_applications" ("article_id");
-- Create index "idx_article_applications_deleted_at" to table: "article_applications"
CREATE INDEX "idx_article_applications_deleted_at" ON "public"."article_applications" ("deleted_at");
-- Create index "idx_article_applications_journal_id" to table: "article_applications"
CREATE INDEX "idx_article_applications_journal_id" ON "public"."article_applications" ("journal_id");
-- Create index "idx_article_applications_number" to table: "article_applications"
CREATE UNIQUE INDEX "idx_article_applications_number" ON "public"."article_applications" ("number");
-- Create index "idx_article_applications_user_id" to table: "article_applications"
CREATE INDEX "idx_article_applications_user_id" ON "public"."article_applications" ("user_id");
-- Create "review_stages" table
CREATE TABLE "public"."review_stages" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "application_id" bigint NULL,
  "stage" integer NOT NULL,
  "status" bigint NOT NULL DEFAULT 0,
  "reason" text NULL,
  "reviewer_id" bigint NULL,
  "reviewed_at" timestamptz NULL,
  "deadline" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "resubmit_deadline" timestamptz NULL,
  "previous_id" bigint NULL,
  "is_old" boolean NOT NULL DEFAULT false,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_article_applications_review_stages" FOREIGN KEY ("application_id") REFERENCES "public"."article_applications" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_review_stages_previous" FOREIGN KEY ("previous_id") REFERENCES "public"."review_stages" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "fk_review_stages_reviewer" FOREIGN KEY ("reviewer_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create index "idx_review_stages_application_id" to table: "review_stages"
CREATE INDEX "idx_review_stages_application_id" ON "public"."review_stages" ("application_id");
-- Create index "idx_review_stages_deleted_at" to table: "review_stages"
CREATE INDEX "idx_review_stages_deleted_at" ON "public"."review_stages" ("deleted_at");
-- Create "anti_plag_result_models" table
CREATE TABLE "public"."anti_plag_result_models" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "review_stage_id" bigint NULL,
  "application_id" bigint NULL,
  "article_id" bigint NULL,
  "journal_id" bigint NULL,
  "external_id" bigint NOT NULL,
  "status" bigint NOT NULL,
  "status_display" character varying(128) NULL,
  "plagiarism_percent" numeric NOT NULL,
  "legal_percent" numeric NOT NULL,
  "self_cite_percent" numeric NOT NULL,
  "unknown_percent" numeric NOT NULL,
  "short_report_url" character varying(2048) NULL,
  "full_report_url" character varying(2048) NULL,
  "external_created_at" timestamptz NULL,
  "certificate" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_anti_plag_result_models_article" FOREIGN KEY ("article_id") REFERENCES "public"."articles" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "fk_anti_plag_result_models_journal" FOREIGN KEY ("journal_id") REFERENCES "public"."journals" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "fk_article_applications_anti_plag_results" FOREIGN KEY ("application_id") REFERENCES "public"."article_applications" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_review_stages_anti_plag_results" FOREIGN KEY ("review_stage_id") REFERENCES "public"."review_stages" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_anti_plag_result_models_deleted_at" to table: "anti_plag_result_models"
CREATE INDEX "idx_anti_plag_result_models_deleted_at" ON "public"."anti_plag_result_models" ("deleted_at");
-- Create index "idx_anti_plag_result_models_external_id" to table: "anti_plag_result_models"
CREATE UNIQUE INDEX "idx_anti_plag_result_models_external_id" ON "public"."anti_plag_result_models" ("external_id");
-- Create "soato_classifiers" table
CREATE TABLE "public"."soato_classifiers" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "code" bigint NULL,
  "name_uz" character varying(255) NULL,
  "name_ru" character varying(255) NULL,
  "name_en" character varying(255) NULL,
  "parent_id" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_soato_classifiers_children" FOREIGN KEY ("parent_id") REFERENCES "public"."soato_classifiers" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_soato_classifiers_code" to table: "soato_classifiers"
CREATE UNIQUE INDEX "idx_soato_classifiers_code" ON "public"."soato_classifiers" ("code");
-- Create index "idx_soato_classifiers_deleted_at" to table: "soato_classifiers"
CREATE INDEX "idx_soato_classifiers_deleted_at" ON "public"."soato_classifiers" ("deleted_at");
-- Create index "idx_soato_classifiers_parent_id" to table: "soato_classifiers"
CREATE INDEX "idx_soato_classifiers_parent_id" ON "public"."soato_classifiers" ("parent_id");
-- Create "organizations" table
CREATE TABLE "public"."organizations" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "tin" character varying(64) NOT NULL,
  "name" character varying(500) NOT NULL,
  "address" character varying(500) NULL,
  "soato_id" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_organizations_soato" FOREIGN KEY ("soato_id") REFERENCES "public"."soato_classifiers" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create index "idx_organizations_deleted_at" to table: "organizations"
CREATE INDEX "idx_organizations_deleted_at" ON "public"."organizations" ("deleted_at");
-- Create index "idx_organizations_soato_id" to table: "organizations"
CREATE INDEX "idx_organizations_soato_id" ON "public"."organizations" ("soato_id");
-- Create "authors" table
CREATE TABLE "public"."authors" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "science_id" character varying(32) NOT NULL,
  "full_name" character varying(255) NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_authors_deleted_at" to table: "authors"
CREATE INDEX "idx_authors_deleted_at" ON "public"."authors" ("deleted_at");
-- Create index "idx_authors_science_id" to table: "authors"
CREATE UNIQUE INDEX "idx_authors_science_id" ON "public"."authors" ("science_id");
-- Create "article_author_affiliations" table
CREATE TABLE "public"."article_author_affiliations" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "article_id" bigint NULL,
  "author_id" bigint NOT NULL,
  "organization_id" bigint NULL,
  "organization_name" character varying(2048) NULL,
  "organization_tin" character varying(32) NULL,
  "position_name" character varying(255) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_article_author_affiliations_organization" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "fk_articles_article_author_affiliations" FOREIGN KEY ("article_id") REFERENCES "public"."articles" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_authors_article_author_affiliations" FOREIGN KEY ("author_id") REFERENCES "public"."authors" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_article_author_affiliations_deleted_at" to table: "article_author_affiliations"
CREATE INDEX "idx_article_author_affiliations_deleted_at" ON "public"."article_author_affiliations" ("deleted_at");
-- Create index "idx_article_id_author_id" to table: "article_author_affiliations"
CREATE UNIQUE INDEX "idx_article_id_author_id" ON "public"."article_author_affiliations" ("article_id", "author_id");
-- Create "article_co_authors" table
CREATE TABLE "public"."article_co_authors" (
  "article_model_id" bigint NOT NULL,
  "author_model_id" bigint NOT NULL,
  PRIMARY KEY ("article_model_id", "author_model_id"),
  CONSTRAINT "fk_article_co_authors_article_model" FOREIGN KEY ("article_model_id") REFERENCES "public"."articles" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_article_co_authors_author_model" FOREIGN KEY ("author_model_id") REFERENCES "public"."authors" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "invoices" table
CREATE TABLE "public"."invoices" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "amount" bigint NOT NULL,
  "payment_provider" character varying(32) NOT NULL,
  "is_paid" boolean NOT NULL DEFAULT false,
  "payment_category" character varying(32) NOT NULL,
  "article_id" bigint NULL,
  "application_id" bigint NULL,
  "user_id" bigint NOT NULL,
  "journal_id" bigint NOT NULL,
  "ticket_file" character varying(256) NULL,
  "ticket_data" jsonb NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_article_applications_invoices" FOREIGN KEY ("application_id") REFERENCES "public"."article_applications" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_invoices_article" FOREIGN KEY ("article_id") REFERENCES "public"."articles" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "fk_invoices_journal" FOREIGN KEY ("journal_id") REFERENCES "public"."journals" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "fk_invoices_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create index "idx_invoices_deleted_at" to table: "invoices"
CREATE INDEX "idx_invoices_deleted_at" ON "public"."invoices" ("deleted_at");
-- Create "article_purchases" table
CREATE TABLE "public"."article_purchases" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "journal_id" bigint NULL,
  "article_id" bigint NULL,
  "user_id" bigint NOT NULL,
  "invoice_id" bigint NOT NULL,
  "amount" bigint NOT NULL DEFAULT 0,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_article_purchases_article" FOREIGN KEY ("article_id") REFERENCES "public"."articles" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_article_purchases_invoice" FOREIGN KEY ("invoice_id") REFERENCES "public"."invoices" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_article_purchases_journal" FOREIGN KEY ("journal_id") REFERENCES "public"."journals" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "fk_article_purchases_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_article_purchase" to table: "article_purchases"
CREATE UNIQUE INDEX "idx_article_purchase" ON "public"."article_purchases" ("article_id", "user_id");
-- Create index "idx_article_purchases_deleted_at" to table: "article_purchases"
CREATE INDEX "idx_article_purchases_deleted_at" ON "public"."article_purchases" ("deleted_at");
-- Create index "idx_article_purchases_journal_id" to table: "article_purchases"
CREATE INDEX "idx_article_purchases_journal_id" ON "public"."article_purchases" ("journal_id");
-- Create "article_requirements" table
CREATE TABLE "public"."article_requirements" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "journal_id" bigint NULL,
  "requirement" character varying(32) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_journals_requirements" FOREIGN KEY ("journal_id") REFERENCES "public"."journals" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create index "idx_article_requirements_deleted_at" to table: "article_requirements"
CREATE INDEX "idx_article_requirements_deleted_at" ON "public"."article_requirements" ("deleted_at");
-- Create index "idx_article_requirements_journal_id" to table: "article_requirements"
CREATE INDEX "idx_article_requirements_journal_id" ON "public"."article_requirements" ("journal_id");
-- Create "study_fields" table
CREATE TABLE "public"."study_fields" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" jsonb NULL,
  "code" bigint NULL,
  "parent_id" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_study_fields_children" FOREIGN KEY ("parent_id") REFERENCES "public"."study_fields" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_study_fields_code" to table: "study_fields"
CREATE UNIQUE INDEX "idx_study_fields_code" ON "public"."study_fields" ("code");
-- Create index "idx_study_fields_deleted_at" to table: "study_fields"
CREATE INDEX "idx_study_fields_deleted_at" ON "public"."study_fields" ("deleted_at");
-- Create index "idx_study_fields_name" to table: "study_fields"
CREATE INDEX "idx_study_fields_name" ON "public"."study_fields" ("name");
-- Create index "idx_study_fields_parent_id" to table: "study_fields"
CREATE INDEX "idx_study_fields_parent_id" ON "public"."study_fields" ("parent_id");
-- Create "article_study_fields" table
CREATE TABLE "public"."article_study_fields" (
  "article_model_id" bigint NOT NULL,
  "study_field_model_id" bigint NOT NULL,
  PRIMARY KEY ("article_model_id", "study_field_model_id"),
  CONSTRAINT "fk_article_study_fields_article_model" FOREIGN KEY ("article_model_id") REFERENCES "public"."articles" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_article_study_fields_study_field_model" FOREIGN KEY ("study_field_model_id") REFERENCES "public"."study_fields" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "tags" table
CREATE TABLE "public"."tags" (
  "id" bigserial NOT NULL,
  "name" character varying(32) NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_tags_name" to table: "tags"
CREATE UNIQUE INDEX "idx_tags_name" ON "public"."tags" ("name");
-- Create "article_tags" table
CREATE TABLE "public"."article_tags" (
  "article_model_id" bigint NOT NULL,
  "tag_model_id" bigint NOT NULL,
  PRIMARY KEY ("article_model_id", "tag_model_id"),
  CONSTRAINT "fk_article_tags_article_model" FOREIGN KEY ("article_model_id") REFERENCES "public"."articles" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_article_tags_tag_model" FOREIGN KEY ("tag_model_id") REFERENCES "public"."tags" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "authorship_claims" table
CREATE TABLE "public"."authorship_claims" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "sender_id" bigint NOT NULL,
  "article_id" bigint NOT NULL,
  "comment" text NULL,
  "status" character varying(20) NULL DEFAULT 'pending',
  "reject_reason" text NULL,
  "reviewed_by_id" bigint NULL,
  "reviewed_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_authorship_claims_article" FOREIGN KEY ("article_id") REFERENCES "public"."articles" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_authorship_claims_reviewed_by" FOREIGN KEY ("reviewed_by_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_authorship_claims_sender" FOREIGN KEY ("sender_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_authorship_claims_deleted_at" to table: "authorship_claims"
CREATE INDEX "idx_authorship_claims_deleted_at" ON "public"."authorship_claims" ("deleted_at");
-- Create index "idx_one_pending_claim" to table: "authorship_claims"
CREATE UNIQUE INDEX "idx_one_pending_claim" ON "public"."authorship_claims" ("sender_id", "article_id") WHERE (((status)::text = 'pending'::text) AND (deleted_at IS NULL));
-- Create "comment_models" table
CREATE TABLE "public"."comment_models" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "article_id" bigint NULL,
  "user_id" bigint NULL,
  "content" text NULL,
  "rating" smallint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_comment_models_article" FOREIGN KEY ("article_id") REFERENCES "public"."articles" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "fk_comment_models_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "chk_comment_models_rating" CHECK ((rating >= 1) AND (rating <= 5))
);
-- Create index "idx_comment_models_deleted_at" to table: "comment_models"
CREATE INDEX "idx_comment_models_deleted_at" ON "public"."comment_models" ("deleted_at");
-- Create "degrees" table
CREATE TABLE "public"."degrees" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "degree_type_id" bigint NOT NULL,
  "degree_type" character varying(255) NULL,
  "field" character varying(64) NULL,
  "degree_status_id" bigint NULL,
  "degree_status_name" character varying(255) NULL,
  "confirmed_date" date NULL,
  "protocol" character varying(64) NULL,
  "user_id" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_degrees_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_degrees_degree_status_id" to table: "degrees"
CREATE INDEX "idx_degrees_degree_status_id" ON "public"."degrees" ("degree_status_id");
-- Create index "idx_degrees_degree_type_id" to table: "degrees"
CREATE INDEX "idx_degrees_degree_type_id" ON "public"."degrees" ("degree_type_id");
-- Create index "idx_degrees_deleted_at" to table: "degrees"
CREATE INDEX "idx_degrees_deleted_at" ON "public"."degrees" ("deleted_at");
-- Create index "idx_user_id_degree_type_id" to table: "degrees"
CREATE UNIQUE INDEX "idx_user_id_degree_type_id" ON "public"."degrees" ("user_id");
-- Create "doctorates" table
CREATE TABLE "public"."doctorates" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "external_id" bigint NOT NULL,
  "dc_type" character varying(64) NULL,
  "edu_lang" character varying(32) NOT NULL,
  "status" character varying(64) NULL,
  "status_code" bigint NULL,
  "admission_year" bigint NULL,
  "direction_name" character varying(255) NOT NULL,
  "direction_code" character varying(16) NOT NULL,
  "advisor_full_name" character varying(255) NULL,
  "advisor_pin" character varying(32) NULL,
  "scientific_work_name" character varying(2048) NULL,
  "organization_tin" character varying(255) NOT NULL,
  "organization_id" bigint NULL,
  "user_id" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_doctorates_organization" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "fk_doctorates_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_doctorates_deleted_at" to table: "doctorates"
CREATE INDEX "idx_doctorates_deleted_at" ON "public"."doctorates" ("deleted_at");
-- Create index "idx_doctorates_external_id" to table: "doctorates"
CREATE UNIQUE INDEX "idx_doctorates_external_id" ON "public"."doctorates" ("external_id");
-- Create index "idx_doctorates_organization_id" to table: "doctorates"
CREATE INDEX "idx_doctorates_organization_id" ON "public"."doctorates" ("organization_id");
-- Create index "idx_doctorates_user_id" to table: "doctorates"
CREATE INDEX "idx_doctorates_user_id" ON "public"."doctorates" ("user_id");
-- Create "drafts" table
CREATE TABLE "public"."drafts" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" bigint NULL,
  "key" character varying(1024) NOT NULL,
  "data" jsonb NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_drafts_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE CASCADE ON DELETE SET NULL
);
-- Create index "idx_drafts_deleted_at" to table: "drafts"
CREATE INDEX "idx_drafts_deleted_at" ON "public"."drafts" ("deleted_at");
-- Create index "idx_drafts_key" to table: "drafts"
CREATE UNIQUE INDEX "idx_drafts_key" ON "public"."drafts" ("key");
-- Create "followers" table
CREATE TABLE "public"."followers" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "following_id" bigint NOT NULL,
  "followed_id" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_followers_followed" FOREIGN KEY ("followed_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_followers_following" FOREIGN KEY ("following_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_followers_deleted_at" to table: "followers"
CREATE INDEX "idx_followers_deleted_at" ON "public"."followers" ("deleted_at");
-- Create index "idx_followers_followed_id" to table: "followers"
CREATE INDEX "idx_followers_followed_id" ON "public"."followers" ("followed_id");
-- Create index "idx_followers_following_id" to table: "followers"
CREATE INDEX "idx_followers_following_id" ON "public"."followers" ("following_id");
-- Create "jobs" table
CREATE TABLE "public"."jobs" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "organization_tin" character varying(32) NULL,
  "organization_name" character varying(2048) NULL,
  "organization_id" bigint NULL,
  "user_id" bigint NOT NULL,
  "position_name" character varying(100) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_jobs_organization" FOREIGN KEY ("organization_id") REFERENCES "public"."organizations" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "fk_jobs_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_jobs_deleted_at" to table: "jobs"
CREATE INDEX "idx_jobs_deleted_at" ON "public"."jobs" ("deleted_at");
-- Create index "idx_jobs_organization_id" to table: "jobs"
CREATE INDEX "idx_jobs_organization_id" ON "public"."jobs" ("organization_id");
-- Create index "idx_jobs_organization_name" to table: "jobs"
CREATE INDEX "idx_jobs_organization_name" ON "public"."jobs" ("organization_name");
-- Create index "idx_jobs_organization_tin" to table: "jobs"
CREATE INDEX "idx_jobs_organization_tin" ON "public"."jobs" ("organization_tin");
-- Create index "idx_jobs_user_id" to table: "jobs"
CREATE INDEX "idx_jobs_user_id" ON "public"."jobs" ("user_id");
-- Create "journal_applications" table
CREATE TABLE "public"."journal_applications" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "journal_id" bigint NULL,
  "user_id" bigint NULL,
  "status" integer NULL DEFAULT 0,
  "rejection_reason" character varying(2048) NULL,
  "reviewed_at" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_journal_applications_journal" FOREIGN KEY ("journal_id") REFERENCES "public"."journals" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "fk_journal_applications_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create index "idx_journal_applications_deleted_at" to table: "journal_applications"
CREATE INDEX "idx_journal_applications_deleted_at" ON "public"."journal_applications" ("deleted_at");
-- Create index "idx_journal_applications_journal_id" to table: "journal_applications"
CREATE INDEX "idx_journal_applications_journal_id" ON "public"."journal_applications" ("journal_id");
-- Create index "idx_journal_applications_user_id" to table: "journal_applications"
CREATE INDEX "idx_journal_applications_user_id" ON "public"."journal_applications" ("user_id");
-- Create "journal_configs" table
CREATE TABLE "public"."journal_configs" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "journal_id" bigint NULL,
  "website_url" text NULL,
  "conf" jsonb NULL,
  "is_active" boolean NOT NULL DEFAULT true,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_journal_configs_journal" FOREIGN KEY ("journal_id") REFERENCES "public"."journals" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_journal_configs_deleted_at" to table: "journal_configs"
CREATE INDEX "idx_journal_configs_deleted_at" ON "public"."journal_configs" ("deleted_at");
-- Create index "idx_journal_configs_journal_id" to table: "journal_configs"
CREATE UNIQUE INDEX "idx_journal_configs_journal_id" ON "public"."journal_configs" ("journal_id");
-- Create index "idx_journal_configs_website_url" to table: "journal_configs"
CREATE UNIQUE INDEX "idx_journal_configs_website_url" ON "public"."journal_configs" ("website_url");
-- Create "journal_indexing" table
CREATE TABLE "public"."journal_indexing" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "indexing_type" character varying(32) NOT NULL,
  "url" character varying(255) NOT NULL,
  "journal_id" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_journals_indexes" FOREIGN KEY ("journal_id") REFERENCES "public"."journals" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_journal_indexing" to table: "journal_indexing"
CREATE INDEX "idx_journal_indexing" ON "public"."journal_indexing" ("indexing_type", "journal_id");
-- Create index "idx_journal_indexing_deleted_at" to table: "journal_indexing"
CREATE INDEX "idx_journal_indexing_deleted_at" ON "public"."journal_indexing" ("deleted_at");
-- Create "journal_many2many_study_fields" table
CREATE TABLE "public"."journal_many2many_study_fields" (
  "journal_model_id" bigint NOT NULL,
  "study_field_model_id" bigint NOT NULL,
  PRIMARY KEY ("journal_model_id", "study_field_model_id"),
  CONSTRAINT "fk_journal_many2many_study_fields_journal_model" FOREIGN KEY ("journal_model_id") REFERENCES "public"."journals" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_journal_many2many_study_fields_study_field_model" FOREIGN KEY ("study_field_model_id") REFERENCES "public"."study_fields" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "journal_ratings" table
CREATE TABLE "public"."journal_ratings" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" bigint NOT NULL,
  "journal_id" bigint NOT NULL,
  "stars" smallint NOT NULL,
  "review" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_journal_ratings_journal" FOREIGN KEY ("journal_id") REFERENCES "public"."journals" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_journal_ratings_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "chk_journal_ratings_stars" CHECK ((stars >= 1) AND (stars <= 5))
);
-- Create index "idx_journal_ratings_deleted_at" to table: "journal_ratings"
CREATE INDEX "idx_journal_ratings_deleted_at" ON "public"."journal_ratings" ("deleted_at");
-- Create index "idx_journal_ratings_journal_id" to table: "journal_ratings"
CREATE INDEX "idx_journal_ratings_journal_id" ON "public"."journal_ratings" ("journal_id");
-- Create index "idx_journal_ratings_user_id" to table: "journal_ratings"
CREATE INDEX "idx_journal_ratings_user_id" ON "public"."journal_ratings" ("user_id");
-- Create "reviewers" table
CREATE TABLE "public"."reviewers" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "external_id" bigint NOT NULL,
  "science_id" text NOT NULL,
  "full_name" character varying(256) NOT NULL,
  "phone_number" character varying(32) NULL,
  "subjects" jsonb NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_reviewers_deleted_at" to table: "reviewers"
CREATE INDEX "idx_reviewers_deleted_at" ON "public"."reviewers" ("deleted_at");
-- Create index "idx_reviewers_external_id" to table: "reviewers"
CREATE UNIQUE INDEX "idx_reviewers_external_id" ON "public"."reviewers" ("external_id");
-- Create index "idx_reviewers_science_id" to table: "reviewers"
CREATE UNIQUE INDEX "idx_reviewers_science_id" ON "public"."reviewers" ("science_id");
-- Create "journal_reviewers" table
CREATE TABLE "public"."journal_reviewers" (
  "journal_model_id" bigint NOT NULL,
  "reviewer_model_id" bigint NOT NULL,
  PRIMARY KEY ("journal_model_id", "reviewer_model_id"),
  CONSTRAINT "fk_journal_reviewers_journal_model" FOREIGN KEY ("journal_model_id") REFERENCES "public"."journals" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_journal_reviewers_reviewer_model" FOREIGN KEY ("reviewer_model_id") REFERENCES "public"."reviewers" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
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
-- Create "news_categories" table
CREATE TABLE "public"."news_categories" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" character varying(255) NOT NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_news_categories_deleted_at" to table: "news_categories"
CREATE INDEX "idx_news_categories_deleted_at" ON "public"."news_categories" ("deleted_at");
-- Create index "idx_news_categories_name" to table: "news_categories"
CREATE UNIQUE INDEX "idx_news_categories_name" ON "public"."news_categories" ("name");
-- Create "news" table
CREATE TABLE "public"."news" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "title" jsonb NOT NULL,
  "body" jsonb NOT NULL,
  "category_id" bigint NOT NULL,
  "image" text NOT NULL,
  "views_count" bigint NOT NULL DEFAULT 0,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_news_category" FOREIGN KEY ("category_id") REFERENCES "public"."news_categories" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create index "idx_news_deleted_at" to table: "news"
CREATE INDEX "idx_news_deleted_at" ON "public"."news" ("deleted_at");
-- Create "notification_models" table
CREATE TABLE "public"."notification_models" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" bigint NULL,
  "title" jsonb NULL,
  "body" jsonb NULL,
  "topic" character varying(32) NULL,
  "extra_data" jsonb NULL,
  "is_email" boolean NOT NULL DEFAULT false,
  "is_sms" boolean NOT NULL DEFAULT false,
  "is_broadcast" boolean NOT NULL DEFAULT false,
  "is_user_read" boolean NOT NULL DEFAULT false,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_notification_models_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_notification_models_deleted_at" to table: "notification_models"
CREATE INDEX "idx_notification_models_deleted_at" ON "public"."notification_models" ("deleted_at");
-- Create index "idx_notification_models_user_id" to table: "notification_models"
CREATE INDEX "idx_notification_models_user_id" ON "public"."notification_models" ("user_id");
-- Create "notification_read_models" table
CREATE TABLE "public"."notification_read_models" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "notification_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_notification_read_models_notification" FOREIGN KEY ("notification_id") REFERENCES "public"."notification_models" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_notification_read_models_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_notification_id_user_id" to table: "notification_read_models"
CREATE UNIQUE INDEX "idx_notification_id_user_id" ON "public"."notification_read_models" ("notification_id", "user_id");
-- Create index "idx_notification_read_models_deleted_at" to table: "notification_read_models"
CREATE INDEX "idx_notification_read_models_deleted_at" ON "public"."notification_read_models" ("deleted_at");
-- Create "notification_token_models" table
CREATE TABLE "public"."notification_token_models" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" bigint NOT NULL,
  "token" character varying(500) NOT NULL,
  "segment" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_notification_token_models_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_notification_token_models_deleted_at" to table: "notification_token_models"
CREATE INDEX "idx_notification_token_models_deleted_at" ON "public"."notification_token_models" ("deleted_at");
-- Create index "idx_notification_token_models_token" to table: "notification_token_models"
CREATE UNIQUE INDEX "idx_notification_token_models_token" ON "public"."notification_token_models" ("token");
-- Create index "idx_notification_token_models_user_id" to table: "notification_token_models"
CREATE INDEX "idx_notification_token_models_user_id" ON "public"."notification_token_models" ("user_id");
-- Create "peer_review_submissions" table
CREATE TABLE "public"."peer_review_submissions" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "external_idempotency_key" text NULL,
  "external_id" bigint NOT NULL,
  "reviewer_internal_id" bigint NULL,
  "reviewer_external_id" bigint NULL,
  "application_id" bigint NULL,
  "sender_title" character varying(1024) NOT NULL,
  "title" character varying(4096) NOT NULL,
  "extra_data" jsonb NULL,
  "old_deadline" timestamptz NULL,
  "deadline" timestamptz NOT NULL,
  "status" bigint NOT NULL DEFAULT 0,
  "review_method" bigint NULL,
  "sender_id" bigint NULL,
  "version" bigint NOT NULL DEFAULT 0,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_peer_review_submissions_application" FOREIGN KEY ("application_id") REFERENCES "public"."article_applications" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "fk_peer_review_submissions_reviewer" FOREIGN KEY ("reviewer_internal_id") REFERENCES "public"."reviewers" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "fk_peer_review_submissions_sender" FOREIGN KEY ("sender_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create index "idx_peer_review_submissions_deleted_at" to table: "peer_review_submissions"
CREATE INDEX "idx_peer_review_submissions_deleted_at" ON "public"."peer_review_submissions" ("deleted_at");
-- Create index "idx_peer_review_submissions_external_id" to table: "peer_review_submissions"
CREATE UNIQUE INDEX "idx_peer_review_submissions_external_id" ON "public"."peer_review_submissions" ("external_id");
-- Create "reference_models" table
CREATE TABLE "public"."reference_models" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" text NULL,
  "article_id" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_reference_models_article" FOREIGN KEY ("article_id") REFERENCES "public"."articles" ("id") ON UPDATE CASCADE ON DELETE SET NULL
);
-- Create index "idx_reference_models_deleted_at" to table: "reference_models"
CREATE INDEX "idx_reference_models_deleted_at" ON "public"."reference_models" ("deleted_at");
-- Create "report" table
CREATE TABLE "public"."report" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "reason" character varying(512) NOT NULL,
  "target_id" bigint NOT NULL,
  "target_type" text NOT NULL,
  "reporter_id" bigint NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_report_reporter" FOREIGN KEY ("reporter_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_report_deleted_at" to table: "report"
CREATE INDEX "idx_report_deleted_at" ON "public"."report" ("deleted_at");
-- Create "research_metrics" table
CREATE TABLE "public"."research_metrics" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" bigint NOT NULL,
  "profile_url" character varying(255) NULL,
  "h_index" bigint NOT NULL DEFAULT 0,
  "source" character varying(32) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_research_metrics_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_research_metrics_deleted_at" to table: "research_metrics"
CREATE INDEX "idx_research_metrics_deleted_at" ON "public"."research_metrics" ("deleted_at");
-- Create index "idx_user_source" to table: "research_metrics"
CREATE UNIQUE INDEX "idx_user_source" ON "public"."research_metrics" ("user_id", "source");
-- Create "roles" table
CREATE TABLE "public"."roles" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" bigint NOT NULL,
  "publisher_id" bigint NULL,
  "journal_id" bigint NULL,
  "role" bigint NULL DEFAULT 0,
  "url" character varying(128) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_roles_journal" FOREIGN KEY ("journal_id") REFERENCES "public"."journals" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_roles_publisher" FOREIGN KEY ("publisher_id") REFERENCES "public"."publishers" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_users_roles" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create index "idx_roles_deleted_at" to table: "roles"
CREATE INDEX "idx_roles_deleted_at" ON "public"."roles" ("deleted_at");
-- Create index "idx_roles_journal_id" to table: "roles"
CREATE INDEX "idx_roles_journal_id" ON "public"."roles" ("journal_id");
-- Create index "idx_roles_publisher_id" to table: "roles"
CREATE INDEX "idx_roles_publisher_id" ON "public"."roles" ("publisher_id");
-- Create index "idx_roles_user_id" to table: "roles"
CREATE INDEX "idx_roles_user_id" ON "public"."roles" ("user_id");
-- Create "spellcheck_results" table
CREATE TABLE "public"."spellcheck_results" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "review_stage_id" bigint NULL,
  "application_id" bigint NULL,
  "journal_id" bigint NULL,
  "file" character varying(255) NOT NULL,
  "result_file" character varying(255) NULL,
  "status" bigint NOT NULL DEFAULT 0,
  "submitter_id" bigint NULL,
  "result_time" timestamptz NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_article_applications_spell_check_results" FOREIGN KEY ("application_id") REFERENCES "public"."article_applications" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_review_stages_spell_check_results" FOREIGN KEY ("review_stage_id") REFERENCES "public"."review_stages" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_spellcheck_results_journal" FOREIGN KEY ("journal_id") REFERENCES "public"."journals" ("id") ON UPDATE NO ACTION ON DELETE SET NULL,
  CONSTRAINT "fk_spellcheck_results_submitter" FOREIGN KEY ("submitter_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL
);
-- Create index "idx_spellcheck_results_application_id" to table: "spellcheck_results"
CREATE INDEX "idx_spellcheck_results_application_id" ON "public"."spellcheck_results" ("application_id");
-- Create index "idx_spellcheck_results_deleted_at" to table: "spellcheck_results"
CREATE INDEX "idx_spellcheck_results_deleted_at" ON "public"."spellcheck_results" ("deleted_at");
-- Create index "idx_spellcheck_results_submitter_id" to table: "spellcheck_results"
CREATE INDEX "idx_spellcheck_results_submitter_id" ON "public"."spellcheck_results" ("submitter_id");
-- Create "support_dialogs" table
CREATE TABLE "public"."support_dialogs" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "message_type" smallint NOT NULL,
  "owner_id" bigint NOT NULL,
  "message" text NOT NULL,
  "chat_id" bigint NOT NULL,
  "is_read" boolean NOT NULL DEFAULT false,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_support_dialogs_owner" FOREIGN KEY ("owner_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_support_dialogs_deleted_at" to table: "support_dialogs"
CREATE INDEX "idx_support_dialogs_deleted_at" ON "public"."support_dialogs" ("deleted_at");
-- Create "transactions" table
CREATE TABLE "public"."transactions" (
  "id" character varying(64) NOT NULL,
  "transaction" character varying(255) NOT NULL,
  "invoice_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "state" bigint NULL,
  "reason" bigint NULL,
  "time" bigint NULL,
  "create_time" bigint NULL,
  "perform_time" bigint NULL,
  "cancel_time" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_transactions_invoice" FOREIGN KEY ("invoice_id") REFERENCES "public"."invoices" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_transactions_invoice_id" to table: "transactions"
CREATE UNIQUE INDEX "idx_transactions_invoice_id" ON "public"."transactions" ("invoice_id");
-- Create "user_profiles" table
CREATE TABLE "public"."user_profiles" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_id" bigint NOT NULL,
  "bio" text NULL,
  "email" character varying(500) NULL,
  "photo" character varying(500) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_user_profiles_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_user_profiles_deleted_at" to table: "user_profiles"
CREATE INDEX "idx_user_profiles_deleted_at" ON "public"."user_profiles" ("deleted_at");
-- Create index "idx_user_profiles_user_id" to table: "user_profiles"
CREATE UNIQUE INDEX "idx_user_profiles_user_id" ON "public"."user_profiles" ("user_id");
-- Create "socials" table
CREATE TABLE "public"."socials" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" character varying(255) NULL,
  "icon" character varying(512) NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_socials_deleted_at" to table: "socials"
CREATE INDEX "idx_socials_deleted_at" ON "public"."socials" ("deleted_at");
-- Create "user_socials" table
CREATE TABLE "public"."user_socials" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "user_profile_id" bigint NOT NULL,
  "social_id" bigint NOT NULL,
  "link" character varying(255) NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_user_profiles_socials" FOREIGN KEY ("user_profile_id") REFERENCES "public"."user_profiles" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_user_socials_social" FOREIGN KEY ("social_id") REFERENCES "public"."socials" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_user_socials_deleted_at" to table: "user_socials"
CREATE INDEX "idx_user_socials_deleted_at" ON "public"."user_socials" ("deleted_at");
-- Create index "idx_user_socials_social_id" to table: "user_socials"
CREATE INDEX "idx_user_socials_social_id" ON "public"."user_socials" ("social_id");
-- Create index "idx_user_socials_user_profile_id" to table: "user_socials"
CREATE INDEX "idx_user_socials_user_profile_id" ON "public"."user_socials" ("user_profile_id");
