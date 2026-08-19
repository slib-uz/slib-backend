-- Modify "journal_editorials" table
ALTER TABLE "public"."journal_editorials" ADD COLUMN "role_code" bigint NOT NULL DEFAULT 1;
