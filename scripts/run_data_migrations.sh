#!/usr/bin/env bash
cd /ugo/scripts/build
./migrate_cid_to_cids
./migrate_position_to_musicalbum_tracks
./set_musicgroup_members_defaults