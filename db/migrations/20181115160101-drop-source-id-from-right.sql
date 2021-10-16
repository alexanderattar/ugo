
-- +migrate Up
ALTER TABLE "right" DROP COLUMN "source_id";

-- +migrate Down
ALTER TABLE "right" ADD COLUMN "source_id" integer NULL;
