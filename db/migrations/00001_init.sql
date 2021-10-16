-- +migrate Up
-- SQL in this section is executed when the migration is applied.
--
-- Create model AudioObject
--
CREATE TABLE "audioobject" ("id" serial NOT NULL PRIMARY KEY, "created_at" timestamp with time zone NOT NULL, "updated_at" timestamp with time zone NOT NULL, "cid" varchar(255) NULL, "type" varchar(255) NOT NULL, "context" varchar(255) NOT NULL, "content_url" varchar(255) NOT NULL, "encoding_format" varchar(255) NOT NULL);
--
-- Create model ImageObject
--
CREATE TABLE "imageobject" ("id" serial NOT NULL PRIMARY KEY, "created_at" timestamp with time zone NOT NULL, "updated_at" timestamp with time zone NOT NULL, "cid" varchar(255) NULL, "type" varchar(255) NOT NULL, "context" varchar(255) NOT NULL, "content_url" varchar(255) NULL, "encoding_format" varchar(255) NULL);
--
-- Create model MusicAlbum
--
CREATE TABLE "musicalbum" ("id" serial NOT NULL PRIMARY KEY, "created_at" timestamp with time zone NOT NULL, "updated_at" timestamp with time zone NOT NULL, "cid" varchar(255) NULL UNIQUE, "type" varchar(255) NOT NULL, "context" varchar(255) NOT NULL, "name" varchar(255) NOT NULL, "album_production_type" varchar(255) NULL, "album_release_type" varchar(255) NULL);
--
-- Create model MusicComposition
--
CREATE TABLE "musiccomposition" ("id" serial NOT NULL PRIMARY KEY, "created_at" timestamp with time zone NOT NULL, "updated_at" timestamp with time zone NOT NULL, "cid" varchar(255) NULL UNIQUE, "type" varchar(255) NOT NULL, "context" varchar(255) NOT NULL, "name" varchar(255) NOT NULL, "iswc" varchar(15) NULL);
--
-- Create model Right
--
CREATE TABLE "right" ("id" serial NOT NULL PRIMARY KEY, "created_at" timestamp with time zone NOT NULL, "updated_at" timestamp with time zone NOT NULL, "cid" varchar(255) NULL UNIQUE, "type" varchar(255) NOT NULL, "context" varchar(255) NOT NULL, "percentage_shares" numeric(5, 2) NOT NULL, "valid_from" varchar(10) NULL, "valid_through" varchar(10) NULL);
--
-- Create model Copyright
--
CREATE TABLE "copyright" ("id" serial NOT NULL PRIMARY KEY, "created_at" timestamp with time zone NOT NULL, "updated_at" timestamp with time zone NOT NULL, "cid" varchar(255) NULL UNIQUE, "type" varchar(255) NOT NULL, "context" varchar(255) NOT NULL, "valid_from" varchar(10) NULL, "valid_through" varchar(10) NULL);
--
-- Create model MusicGroup
--
CREATE TABLE "musicgroup" ("id" serial NOT NULL PRIMARY KEY, "created_at" timestamp with time zone NOT NULL, "updated_at" timestamp with time zone NOT NULL, "cid" varchar(255) NULL UNIQUE, "type" varchar(255) NOT NULL, "context" varchar(255) NOT NULL, "name" varchar(255) NOT NULL, "description" text NOT NULL, "email" varchar(254) NULL);
--
-- Create model MusicPlaylist
--
CREATE TABLE "musicplaylist" ("id" serial NOT NULL PRIMARY KEY, "created_at" timestamp with time zone NOT NULL, "updated_at" timestamp with time zone NOT NULL, "cid" varchar(255) NULL UNIQUE, "type" varchar(255) NOT NULL, "context" varchar(255) NOT NULL, "name" varchar(255) NOT NULL);
--
-- Create model MusicRecording
--
CREATE TABLE "musicrecording" ("id" serial NOT NULL PRIMARY KEY, "created_at" timestamp with time zone NOT NULL, "updated_at" timestamp with time zone NOT NULL, "cid" varchar(255) NULL UNIQUE, "type" varchar(255) NOT NULL, "context" varchar(255) NOT NULL, "name" varchar(255) NOT NULL, "duration" varchar(255) NULL, "isrc" varchar(12) NULL, "position" varchar(255) NULL, genres text[] NULL DEFAULT '{}'::text[]);
--
-- Create model MusicRelease
--
CREATE TABLE "musicrelease" ("id" serial NOT NULL PRIMARY KEY, "created_at" timestamp with time zone NOT NULL, "updated_at" timestamp with time zone NOT NULL, "cid" varchar(255) NULL UNIQUE, "type" varchar(255) NOT NULL, "context" varchar(255) NOT NULL, "description" text NULL, "date_published" varchar(255) NULL, "catalog_number" varchar(255) NULL, "music_release_format" varchar(255) NULL, "price" numeric(8, 2));
--
-- Create model Purchase
--
CREATE TABLE "purchase" ("id" serial NOT NULL PRIMARY KEY, "created_at" timestamp with time zone NOT NULL, "updated_at" timestamp with time zone NOT NULL, "cid" varchar(255) NULL UNIQUE, "type" varchar(255) NOT NULL, "context" varchar(255) NOT NULL, "tx_hash" varchar(255) NOT NULL, "buyer_id" integer NOT NULL, "musicrelease_id" integer NOT NULL);
--
-- Create model Organization
--
CREATE TABLE "organization" ("id" serial NOT NULL PRIMARY KEY, "created_at" timestamp with time zone NOT NULL, "updated_at" timestamp with time zone NOT NULL, "cid" varchar(255) NULL UNIQUE, "type" varchar(255) NOT NULL, "context" varchar(255) NOT NULL, "name" varchar(255) NOT NULL, "description" text NOT NULL, "email" varchar(254) NULL);
--
-- Create model Person
--
CREATE TABLE "person" ("id" serial NOT NULL PRIMARY KEY, "created_at" timestamp with time zone NOT NULL, "updated_at" timestamp with time zone NOT NULL, "cid" varchar(255) NULL, "type" varchar(255) NOT NULL, "context" varchar(255) NOT NULL, "ethereum_address" varchar(255) NOT NULL UNIQUE, "given_name" varchar(255) NULL, "family_name" varchar(255) NULL, "email" varchar(255) NULL);
--
-- Add field release_of to musicrelease
--
ALTER TABLE "musicgroup" ADD COLUMN "image_id" integer NOT NULL;
--
-- Add field record_label to musicrelease
--
ALTER TABLE "musicrelease" ADD COLUMN "record_label_id" integer NULL UNIQUE;
--
-- Add field release_of to musicrelease
--
ALTER TABLE "musicrelease" ADD COLUMN "release_of_id" integer NOT NULL;
--
-- Add field release_of to musicrelease
--
ALTER TABLE "musicrelease" ADD COLUMN "image_id" integer NOT NULL;
--
-- Add field audio to musicrecording
--
ALTER TABLE "musicrecording" ADD COLUMN "audio_id" integer NOT NULL UNIQUE;
--
-- Add field by_artist to musicrecording
--
ALTER TABLE "musicrecording" ADD COLUMN "by_artist_id" integer NOT NULL;
--
-- Add field recording_of to musicrecording
--
ALTER TABLE "musicrecording" ADD COLUMN "recording_of_id" integer NOT NULL;
--
-- Add field musicrelease to copyright
--
ALTER TABLE "copyright" ADD COLUMN "musicrelease_id" integer NOT NULL;
--
-- Add field source to right
--
ALTER TABLE "right" ADD COLUMN "source_id" integer NOT NULL;
--
-- Add field by_artist to musicalbum
--
ALTER TABLE "musicalbum" ADD COLUMN "by_artist_id" integer NULL;
--
-- Add field tracks to musicplaylist
--
CREATE TABLE "musicplaylist_tracks" ("id" serial NOT NULL PRIMARY KEY, "musicplaylist_id" integer NOT NULL, "musicrecording_id" integer NOT NULL);
--
-- Add field members to musicgroup
--
CREATE TABLE "musicgroup_members" ("id" serial NOT NULL PRIMARY KEY, "musicgroup_id" integer NOT NULL, "person_id" integer NOT NULL);
--
-- Add fields person and organization to right
--
CREATE TABLE "right_party" ("id" serial NOT NULL PRIMARY KEY, "right_id" integer NOT NULL, "person_id" integer NULL, "organization_id" integer NULL);
--
-- Add field composer to musiccomposition
--
CREATE TABLE "musiccomposition_composer" ("id" serial NOT NULL PRIMARY KEY, "musiccomposition_id" integer NOT NULL, "person_id" integer NOT NULL);
--
-- Add field tracks to musicalbum
--
CREATE TABLE "musicalbum_tracks" ("id" serial NOT NULL PRIMARY KEY, "musicalbum_id" integer NOT NULL, "musicrecording_id" integer NOT NULL);
--
-- Create index audioobject_cid_idx on field(s) cid of model audioobject
--
CREATE INDEX "audioobject_cid_idx" ON "audioobject" ("cid");
--
-- Create index musicrelease_cid_idx on field(s) cid of model musicrelease
--
CREATE INDEX "musicrelease_cid_idx" ON "musicrelease" ("cid");
--
-- Create index purchase_cid_idx on field(s) cid of model purchase
--
CREATE INDEX "purchase_cid_idx" ON "purchase" ("cid");
--
-- Create index musicrecording_cid_idx on field(s) cid of model musicrecording
--
CREATE INDEX "musicrecording_cid_idx" ON "musicrecording" ("cid");
--
-- Create index musicplaylist_cid_idx on field(s) cid of model musicplaylist
--
CREATE INDEX "musicplaylist_cid_idx" ON "musicplaylist" ("cid");
--
-- Create index musicgroup_cid_idx on field(s) cid of model musicgroup
--
CREATE INDEX "musicgroup_cid_idx" ON "musicgroup" ("cid");
--
-- Create index right_cid_idx on field(s) cid of model right
--
CREATE INDEX "right_cid_idx" ON "right" ("cid");
--
-- Create index copyright_cid_idx on field(s) cid of model copyright
--
CREATE INDEX "copyright_cid_idx" ON "copyright" ("cid");
--
-- Create index musiccomposition_cid_idx on field(s) cid of model musiccomposition
--
CREATE INDEX "musiccomposition_cid_idx" ON "musiccomposition" ("cid");
--
-- Create index musicalbum_cid_idx on field(s) cid of model musicalbum
--
CREATE INDEX "musicalbum_cid_idx" ON "musicalbum" ("cid");
--
-- Create index imageobject_cid_idx on field(s) cid of model imageobject
--
CREATE INDEX "imageobject_cid_idx" ON "imageobject" ("cid");
--
-- Create index person_cid_idx on field(s) cid of model person
--
CREATE INDEX "person_cid_idx" ON "person" ("cid");
--
-- Create index organization_cid_idx on field(s) cid of model organization
--
CREATE INDEX "organization_cid_idx" ON "organization" ("cid");
CREATE INDEX "audioobject_cid_like" ON "audioobject" ("cid" varchar_pattern_ops);
CREATE INDEX "imageobject_cid_like" ON "imageobject" ("cid" varchar_pattern_ops);
CREATE INDEX "musicalbum_cid_like" ON "musicalbum" ("cid" varchar_pattern_ops);
CREATE INDEX "musiccomposition_cid_like" ON "musiccomposition" ("cid" varchar_pattern_ops);
CREATE INDEX "right_cid_like" ON "right" ("cid" varchar_pattern_ops);
CREATE INDEX "copyright_cid_like" ON "copyright" ("cid" varchar_pattern_ops);
CREATE INDEX "musicgroup_cid_like" ON "musicgroup" ("cid" varchar_pattern_ops);
CREATE INDEX "musicplaylist_cid_like" ON "musicplaylist" ("cid" varchar_pattern_ops);
CREATE INDEX "musicrecording_cid_like" ON "musicrecording" ("cid" varchar_pattern_ops);
CREATE INDEX "musicrelease_cid_like" ON "musicrelease" ("cid" varchar_pattern_ops);
CREATE INDEX "organization_cid_like" ON "organization" ("cid" varchar_pattern_ops);
CREATE INDEX "person_cid_like" ON "person" ("cid" varchar_pattern_ops);
CREATE INDEX "person_ethereum_address_like" ON "person" ("ethereum_address" varchar_pattern_ops);
ALTER TABLE "purchase" ADD FOREIGN KEY ("buyer_id") REFERENCES "person" ("id") DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "purchase" ADD FOREIGN KEY ("musicrelease_id") REFERENCES "musicrelease" ("id") DEFERRABLE INITIALLY DEFERRED;
CREATE INDEX "purchase_cid_like" ON "purchase" ("cid" varchar_pattern_ops);
CREATE INDEX "purchase_buyer_id_idx" ON "purchase" ("buyer_id");
CREATE INDEX "purchase_musicrelease_id_idx" ON "purchase" ("musicrelease_id");
ALTER TABLE "musicrelease" ADD FOREIGN KEY ("record_label_id") REFERENCES "organization" ("id") DEFERRABLE INITIALLY DEFERRED;
CREATE INDEX "musicrelease_release_of_id_idx" ON "musicrelease" ("release_of_id");
ALTER TABLE "musicrelease" ADD FOREIGN KEY ("release_of_id") REFERENCES "musicalbum" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;
CREATE INDEX "musicrelease_image_id_idx" ON "musicrelease" ("image_id");
ALTER TABLE "musicrelease" ADD FOREIGN KEY ("image_id") REFERENCES "imageobject" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "musicrecording" ADD FOREIGN KEY ("audio_id") REFERENCES "audioobject" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;
CREATE INDEX "musicrecording_by_artist_id_idx" ON "musicrecording" ("by_artist_id");
ALTER TABLE "musicrecording" ADD FOREIGN KEY ("by_artist_id") REFERENCES "musicgroup" ("id") DEFERRABLE INITIALLY DEFERRED;
CREATE INDEX "musicrecording_recording_of_id_idx" ON "musicrecording" ("recording_of_id");
ALTER TABLE "musicrecording" ADD FOREIGN KEY ("recording_of_id") REFERENCES "musiccomposition" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "musicplaylist_tracks" ADD FOREIGN KEY ("musicplaylist_id") REFERENCES "musicplaylist" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "musicplaylist_tracks" ADD FOREIGN KEY ("musicrecording_id") REFERENCES "musicrecording" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "musicplaylist_tracks" ADD UNIQUE ("musicplaylist_id", "musicrecording_id");
CREATE INDEX "musicplaylist_tracks_musicplaylist_id_idx" ON "musicplaylist_tracks" ("musicplaylist_id");
CREATE INDEX "musicplaylist_tracks_musicrecording_id_idx" ON "musicplaylist_tracks" ("musicrecording_id");
ALTER TABLE "musicgroup_members" ADD FOREIGN KEY ("musicgroup_id") REFERENCES "musicgroup" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "musicgroup_members" ADD FOREIGN KEY ("person_id") REFERENCES "person" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "musicgroup_members" ADD UNIQUE ("musicgroup_id", "person_id");
CREATE INDEX "musicgroup_members_musicgroup_id_idx" ON "musicgroup_members" ("musicgroup_id");
CREATE INDEX "musicgroup_members_person_id_idx" ON "musicgroup_members" ("person_id");
CREATE INDEX "musicgroup_image_id_idx" ON "musicgroup" ("image_id");
ALTER TABLE "musicgroup" ADD FOREIGN KEY ("image_id") REFERENCES "imageobject" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "right_party" ADD FOREIGN KEY ("right_id") REFERENCES "right" ("id") DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "right_party" ADD FOREIGN KEY ("person_id") REFERENCES "person" ("id") DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "right_party" ADD FOREIGN KEY ("organization_id") REFERENCES "organization" ("id") DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "right_party" ADD UNIQUE ("right_id", "person_id", "organization_id");
CREATE INDEX "right_party_right_id_idx" ON "right_party" ("right_id");
CREATE INDEX "right_party_person_id_idx" ON "right_party" ("person_id");
CREATE INDEX "right_party_organization_id_idx" ON "right_party" ("organization_id");
ALTER TABLE "musiccomposition_composer" ADD FOREIGN KEY ("musiccomposition_id") REFERENCES "musiccomposition" ("id") DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "musiccomposition_composer" ADD FOREIGN KEY ("person_id") REFERENCES "person" ("id") DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "musiccomposition_composer" ADD UNIQUE ("musiccomposition_id", "person_id");
CREATE INDEX "musiccomposition_composer_musiccomposition_id_idx" ON "musiccomposition_composer" ("musiccomposition_id");
CREATE INDEX "musiccomposition_composer_person_id_idx" ON "musiccomposition_composer" ("person_id");
CREATE INDEX "musicalbum_by_artist_id_idx" ON "musicalbum" ("by_artist_id");
ALTER TABLE "musicalbum" ADD FOREIGN KEY ("by_artist_id") REFERENCES "musicgroup" ("id") DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "musicalbum_tracks" ADD FOREIGN KEY ("musicalbum_id") REFERENCES "musicalbum" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "musicalbum_tracks" ADD FOREIGN KEY ("musicrecording_id") REFERENCES "musicrecording" ("id") ON DELETE CASCADE DEFERRABLE INITIALLY DEFERRED;
ALTER TABLE "musicalbum_tracks" ADD UNIQUE ("musicalbum_id", "musicrecording_id");
CREATE INDEX "musicalbum_tracks_musicalbum_id_idx" ON "musicalbum_tracks" ("musicalbum_id");
CREATE INDEX "musicalbum_tracks_musicrecording_id_idx" ON "musicalbum_tracks" ("musicrecording_id");

-- +migrate Down
-- SQL in this section is executed when the migration is rolled back.
--
-- Delete model Person
--
DROP TABLE "person" CASCADE;
--
-- Delete model Organization
--
DROP TABLE "organization" CASCADE;
--
-- Delete model Purchase
--
DROP TABLE "purchase" CASCADE;
--
-- Delete model MusicRelease
--
DROP TABLE "musicrelease" CASCADE;
--
-- Delete model MusicRecording
--
DROP TABLE "musicrecording" CASCADE;
--
-- Delete model MusicPlaylist
--
DROP TABLE "musicplaylist" CASCADE;
--
-- Delete model MusicGroup
--
DROP TABLE "musicgroup" CASCADE;
--
-- Delete model Right
--
DROP TABLE "copyright" CASCADE;
--
-- Delete model Right
--
DROP TABLE "right" CASCADE;
--
-- Delete model MusicComposition
--
DROP TABLE "musiccomposition" CASCADE;
--
-- Delete model MusicAlbum
--
DROP TABLE "musicalbum" CASCADE;
--
-- Delete model ImageObject
--
DROP TABLE "imageobject";
--
-- Delete model AudioObject
--
DROP TABLE "audioobject" CASCADE;

--
-- Delete JOIN tables
--
DROP TABLE "musicalbum_tracks" CASCADE;
DROP TABLE "musiccomposition_composer" CASCADE;
DROP TABLE "right_party" CASCADE;
DROP TABLE "musicgroup_members" CASCADE;
DROP TABLE "musicplaylist_tracks" CASCADE;
