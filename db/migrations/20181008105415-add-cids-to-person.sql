-- +migrate Up
ALTER TABLE person ADD COLUMN cids text[] DEFAULT '{}'::text[];

-- +migrate Down
ALTER TABLE person DROP COLUMN "cids";
