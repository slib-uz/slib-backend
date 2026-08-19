-- Modify "journal_configs" table
ALTER TABLE "public"."journal_configs" ADD COLUMN "creator_id" bigint NULL, ADD CONSTRAINT "fk_journal_configs_creator" FOREIGN KEY ("creator_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE SET NULL;
