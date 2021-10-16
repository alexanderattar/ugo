
-- +migrate Up
--
-- Add field image_id to musicrecording
--
ALTER TABLE "musicrecording" ADD COLUMN "image_id" integer NULL;
CREATE INDEX "musicrecording_image_id_idx" ON "musicrecording" ("image_id");
ALTER TABLE "musicrecording" ADD FOREIGN KEY ("image_id") REFERENCES "imageobject" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;

-- +migrate Down
ALTER TABLE musicrecording DROP COLUMN image_id CASCADE;