-- +migrate Up
ALTER TABLE musicgroup_members ADD COLUMN "percentage_shares" numeric(5, 2) NULL;

-- +migrate Down
ALTER TABLE musicgroup_members DROP COLUMN "percentage_shares";
