-- +migrate Up
ALTER TABLE musicplaylist ADD COLUMN cids text[] DEFAULT '{}'::text[];

-- +migrate Down
ALTER TABLE musicplaylist DROP COLUMN "cids";
