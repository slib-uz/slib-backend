-- Modify "article_purchases" table
ALTER TABLE "public"."article_purchases" DROP COLUMN "invoice_id";
-- Drop "transactions" table
DROP TABLE "public"."transactions";
-- Drop "invoices" table
DROP TABLE "public"."invoices";
