-- Modify "user_profiles" table
ALTER TABLE "public"."user_profiles" ADD COLUMN "last_online_at" timestamptz NULL;
