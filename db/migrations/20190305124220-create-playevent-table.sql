
-- +migrate Up

CREATE TABLE "playevent" ("id" serial NOT NULL PRIMARY KEY, "created_at" timestamp with time zone NOT NULL, "updated_at" timestamp with time zone NOT NULL);
ALTER TABLE "playevent" ADD COLUMN "playedby_id" integer NULL;
ALTER TABLE "playevent" ADD COLUMN "musicrecording_id" integer NOT NULL;

CREATE INDEX "playevent_playedby_id_idx" ON "playevent" ("playedby_id");
ALTER TABLE "playevent" ADD FOREIGN KEY ("playedby_id") REFERENCES "person" ("id") DEFERRABLE INITIALLY DEFERRED;

CREATE INDEX "playevent_musicrecording_id_idx" ON "playevent" ("musicrecording_id");
ALTER TABLE "playevent" ADD FOREIGN KEY ("musicrecording_id") REFERENCES "musicrecording" ("id") DEFERRABLE INITIALLY DEFERRED;

-- +migrate Down

DROP TABLE "playevent" CASCADE;