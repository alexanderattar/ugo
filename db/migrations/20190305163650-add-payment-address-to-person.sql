
-- +migrate Up
ALTER TABLE person ADD COLUMN "payment_address" varchar(255) NULL;

-- +migrate Down
ALTER TABLE person DROP COLUMN "payment_address";
