
-- +migrate Up
ALTER TABLE musicgroup_members ADD COLUMN "description" varchar NULL;

-- +migrate Down
ALTER TABLE musicgroup_members DROP COLUMN "description";