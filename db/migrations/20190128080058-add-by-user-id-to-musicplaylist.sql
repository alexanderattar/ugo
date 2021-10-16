-- +migrate Up
ALTER TABLE "musicplaylist" ADD COLUMN "by_user_id" integer NULL;
ALTER TABLE "musicplaylist" ADD FOREIGN KEY ("by_user_id") REFERENCES "person" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;
CREATE INDEX "musicplaylist_by_user_id_idx" ON "musicplaylist" ("by_user_id");

-- +migrate Down
ALTER TABLE "musicplaylist" DROP COLUMN "by_user_id";
