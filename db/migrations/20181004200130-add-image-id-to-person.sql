
-- +migrate Up
--
-- Add field image_id to person
--
ALTER TABLE "person" ADD COLUMN "image_id" integer NULL;
CREATE INDEX "person_image_id_idx" ON "person" ("image_id");
ALTER TABLE "person" ADD FOREIGN KEY ("image_id") REFERENCES "imageobject" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;

-- +migrate Down
ALTER TABLE person DROP COLUMN image_id CASCADE;
