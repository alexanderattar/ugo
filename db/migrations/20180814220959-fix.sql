
-- +migrate Up

ALTER TABLE "musicrecording" DROP CONSTRAINT musicrecording_audio_id_fkey;
ALTER TABLE "musicrecording" DROP CONSTRAINT musicrecording_image_id_fkey;
ALTER TABLE "musicrecording" DROP CONSTRAINT musicrecording_recording_of_id_fkey;

-- +migrate Down
