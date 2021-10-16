
-- +migrate Up

CREATE TABLE "signedmessage" ("id" serial NOT NULL PRIMARY KEY, "created_at" timestamp with time zone NOT NULL, "updated_at" timestamp with time zone NOT NULL, "message" text NULL, "ethereum_address" varchar(255) NOT NULL, "signature" varchar(255) NOT NULL);

-- +migrate Down

DROP TABLE "signedmessage" CASCADE;
