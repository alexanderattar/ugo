-- +migrate Up
ALTER TABLE musicgroup_members ADD COLUMN "musicgroup_admin" BOOLEAN NULL DEFAULT FALSE;

-- +migrate Down
ALTER TABLE musicgroup_members DROP COLUMN "musicgroup_admin";
