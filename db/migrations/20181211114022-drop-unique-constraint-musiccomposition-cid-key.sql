-- +migrate Up
ALTER TABLE "musiccomposition" DROP CONSTRAINT "musiccomposition_cid_key";

-- +migrate Down
ALTER TABLE "musiccomposition" ADD CONSTRAINT "musiccomposition_cid_key" UNIQUE ("cid");