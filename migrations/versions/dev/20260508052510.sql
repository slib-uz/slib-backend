-- Modify "clients" table
ALTER TABLE "public"."clients" ADD COLUMN "journal_id" bigint NULL, ADD CONSTRAINT "fk_clients_journal" FOREIGN KEY ("journal_id") REFERENCES "public"."journals" ("id") ON UPDATE NO ACTION ON DELETE SET NULL;
-- Create index "idx_clients_journal_id" to table: "clients"
CREATE INDEX "idx_clients_journal_id" ON "public"."clients" ("journal_id");
