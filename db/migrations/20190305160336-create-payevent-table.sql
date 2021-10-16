
-- +migrate Up

CREATE TABLE "payevent" ("id" serial NOT NULL PRIMARY KEY, "created_at" timestamp with time zone NOT NULL, "updated_at" timestamp with time zone NOT NULL, "amount" numeric(22, 18), "link" varchar(255) NULL);

-- Add join columns
ALTER TABLE "payevent" ADD COLUMN "playedby_id" integer NULL;
ALTER TABLE "payevent" ADD COLUMN "beneficiary_id" integer NULL;
ALTER TABLE "payevent" ADD COLUMN "musicrecording_id" integer NULL;

-- Add join indexes and link with foreign keys to each join column
-- CREATE INDEX "payevent_playedby_id_idx" ON "payevent" ("playedby_id");
ALTER TABLE "payevent" ADD FOREIGN KEY ("playedby_id") REFERENCES "person" ("id") DEFERRABLE INITIALLY DEFERRED;

-- CREATE INDEX "payevent_beneficiary_id_idx" ON "payevent" ("beneficiary_id");
ALTER TABLE "payevent" ADD FOREIGN KEY ("beneficiary_id") REFERENCES "person" ("id") DEFERRABLE INITIALLY DEFERRED;

-- CREATE INDEX "payevent_musicrecording_id_idx" ON "payevent" ("musicrecording_id");
ALTER TABLE "payevent" ADD FOREIGN KEY ("musicrecording_id") REFERENCES "musicrecording" ("id") DEFERRABLE INITIALLY DEFERRED;


-- +migrate Down

DROP TABLE "payevent" CASCADE;
