-- Rename a column from "do_iprefix" to "doi_prefix"
ALTER TABLE "public"."journal_doi_settings" RENAME COLUMN "do_iprefix" TO "doi_prefix";
