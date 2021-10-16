
-- +migrate Up

CREATE TABLE "report" ("id" serial NOT NULL PRIMARY KEY, "created_at" timestamp with time zone NOT NULL, "updated_at" timestamp with time zone NOT NULL, "state" varchar(256) NOT NULL, "reason" varchar(256) NOT NULL, "message" varchar(512) NULL, "response" varchar(512) NULL, "email" varchar(256) NULL, "reporter_id" integer NULL, "musicrelease_id" integer NULL);

ALTER TABLE "report" ADD FOREIGN KEY ("reporter_id") REFERENCES "person" ("id") DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "report" ADD FOREIGN KEY ("musicrelease_id") REFERENCES "musicrelease" ("id") DEFERRABLE INITIALLY DEFERRED;

CREATE INDEX "report_reporter_id_idx" ON "report" ("reporter_id");
CREATE INDEX "report_musicrelease_id_idx" ON "report" ("musicrelease_id");

-- +migrate Down

DROP TABLE "report" CASCADE;
