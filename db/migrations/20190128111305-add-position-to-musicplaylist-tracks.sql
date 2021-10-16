-- +migrate Up
ALTER TABLE musicplaylist_tracks ADD COLUMN "position" integer NULL;

-- +migrate Down
ALTER TABLE musicplaylist_tracks DROP COLUMN "position";
