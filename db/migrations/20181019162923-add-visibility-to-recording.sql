
-- +migrate Up
ALTER TABLE musicrecording ADD COLUMN "visibility" varchar NULL;

-- +migrate Down
ALTER TABLE musicrecording DROP COLUMN "visibility";