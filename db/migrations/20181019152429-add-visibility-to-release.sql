
-- +migrate Up
ALTER TABLE musicrelease ADD COLUMN "visibility" varchar NULL;

-- +migrate Down
ALTER TABLE musicrelease DROP COLUMN "visibility";