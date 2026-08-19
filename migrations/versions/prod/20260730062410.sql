-- Modify "journal_editorials" table
ALTER TABLE "public"."journal_editorials" ADD COLUMN "order" bigint NOT NULL DEFAULT 0;
