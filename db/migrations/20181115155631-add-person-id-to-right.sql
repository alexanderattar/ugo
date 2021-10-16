-- +migrate Up
ALTER TABLE "right" ADD COLUMN "person_id" integer NULL;
ALTER TABLE "right" ADD FOREIGN KEY ("person_id") REFERENCES "person" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;
CREATE INDEX "right_person_id_idx" ON "right" ("person_id");

-- +migrate Down
ALTER TABLE "right" DROP COLUMN "person_id";
