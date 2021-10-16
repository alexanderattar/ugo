-- +migrate Up
-- SQL in this section is executed when the migration is applied.
--
-- Add active field
--
ALTER TABLE musicrelease
ADD COLUMN "active" BOOLEAN NOT NULL DEFAULT TRUE;

-- +migrate Down
-- SQL in this section is executed when the migration is rolled back.
--
-- Remove active field
--
ALTER TABLE musicrelease
DROP COLUMN "active";