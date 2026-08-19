-- Modify "doi_deposits" table
ALTER TABLE "public"."doi_deposits" ADD COLUMN "request_body" text NULL, ADD COLUMN "response_body" text NULL;
