
-- +migrate Up

ALTER TABLE "musicrecording" ADD CONSTRAINT scratch FOREIGN KEY (audio_id) REFERENCES audioobject(id);
ALTER TABLE "musicrecording" DROP CONSTRAINT musicrecording_audio_id_fkey;
ALTER TABLE "musicrecording" RENAME CONSTRAINT scratch TO musicrecording_audio_id_fkey;

ALTER TABLE "musicrecording" ADD CONSTRAINT scratch FOREIGN KEY (image_id) REFERENCES imageobject(id);
ALTER TABLE "musicrecording" DROP CONSTRAINT musicrecording_image_id_fkey;
ALTER TABLE "musicrecording" RENAME CONSTRAINT scratch TO musicrecording_image_id_fkey;

ALTER TABLE "musicrecording" ADD CONSTRAINT scratch FOREIGN KEY (recording_of_id) REFERENCES musiccomposition(id);
ALTER TABLE "musicrecording" DROP CONSTRAINT musicrecording_recording_of_id_fkey;
ALTER TABLE "musicrecording" RENAME CONSTRAINT scratch TO musicrecording_recording_of_id_fkey;

-- +migrate Down
