-- Modify "doi_deposits" table
ALTER TABLE "public"."doi_deposits" DROP CONSTRAINT "fk_doi_deposits_article", ALTER COLUMN "article_id" DROP NOT NULL, ADD CONSTRAINT "fk_doi_deposits_article" FOREIGN KEY ("article_id") REFERENCES "public"."articles" ("id") ON UPDATE NO ACTION ON DELETE SET NULL;
