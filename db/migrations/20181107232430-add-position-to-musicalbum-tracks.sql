-- +migrate Up
ALTER TABLE musicalbum_tracks ADD COLUMN "position" integer NULL;

-- +migrate Down
ALTER TABLE musicalbum_tracks DROP COLUMN "position";
