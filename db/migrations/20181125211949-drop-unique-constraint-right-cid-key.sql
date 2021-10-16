
-- +migrate Up
ALTER TABLE "right" DROP CONSTRAINT "right_cid_key";

-- +migrate Down
ALTER TABLE "right" ADD CONSTRAINT "right_cid_key" UNIQUE ("cid");
