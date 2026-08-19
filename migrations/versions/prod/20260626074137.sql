-- Modify "articles" table
ALTER TABLE "public"."articles" ADD COLUMN "unconfirmed_authors" jsonb NULL;
