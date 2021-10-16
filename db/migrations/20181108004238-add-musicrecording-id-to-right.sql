
-- +migrate Up
ALTER TABLE "right" ADD COLUMN "musicrecording_id" integer NULL;
ALTER TABLE "right" ADD FOREIGN KEY ("musicrecording_id") REFERENCES "musicrecording" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;
CREATE INDEX "right_musicrecording_id_idx" ON "right" ("musicrecording_id");

-- +migrate Down
ALTER TABLE "right" DROP COLUMN "musicrecording_id";
