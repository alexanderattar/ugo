
-- +migrate Up

ALTER TABLE "person" DROP CONSTRAINT person_ethereum_address_key;

-- +migrate Down
