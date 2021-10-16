
-- +migrate Up
ALTER TABLE "musicplaylist" ADD COLUMN "image_id" integer NULL;
ALTER TABLE "musicplaylist" ADD FOREIGN KEY ("image_id") REFERENCES "imageobject" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;
CREATE INDEX "musicplaylist_image_id_idx" ON "musicplaylist" ("image_id");

-- +migrate Down
ALTER TABLE "musicplaylist" DROP COLUMN "image_id";
