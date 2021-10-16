-- +migrate Up
ALTER TABLE musicgroup ADD COLUMN cids text[] NOT NULL DEFAULT '{}'::text[];

-- +migrate Down
ALTER TABLE musicgroup DROP COLUMN "cids";