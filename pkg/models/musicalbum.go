package models

// MusicAlbum model
import (
	"errors"
	"fmt"
	"time"

	"github.com/consensys/ugo/pkg/lg"
	"github.com/jmoiron/sqlx"
)

// MusicAlbum model
type MusicAlbum struct {
	ID                  IDType            `json:"id"`
	CID                 string            `json:"cid"`
	Type                string            `json:"@type"`
	Context             string            `json:"@context"`
	CreatedAt           time.Time         `json:"createdAt" db:"created_at"`
	UpdatedAt           time.Time         `json:"updatedAt"  db:"updated_at"`
	Name                *string           `json:"name"`
	Tracks              []*MusicRecording `json:"tracks"`
	AlbumProductionType *string           `json:"albumProductionType" db:"album_production_type"`
	AlbumReleaseType    *string           `json:"albumReleaseType"  db:"album_release_type"`
	ByArtist            *MusicGroup       `json:"byArtist"`
	ByArtistID          *IDType           `json:"-" db:"by_artist_id"`
}

// All gets all of the given MusicRelease objects
func (obj *MusicAlbum) All(db *sqlx.DB) ([]*MusicAlbum, error) {
	return nil, errors.New("Not implemented")
}

// Get MusicAlbum by ID
func (obj *MusicAlbum) Get(db *sqlx.DB, id IDType) (*MusicAlbum, error) {
	query := &SelectQuery{}
	query.SetQuery(`
    SELECT *
    FROM musicalbum
    WHERE musicalbum.id=$1
    `, id)

	ma := &MusicAlbum{}
	err := db.Get(ma, query.Query(), query.Args()...)
	if err != nil {
		return nil, fmt.Errorf("Error getting MusicAlbum (%v)", err)
	}

	ma.Tracks, err = ma.GetTracks(db, ma.ID)
	if err != nil {
		lg.Errorf("Error GetTracks (%v)", err)
	}

	ma.ByArtist, err = ma.GetArtist(db, *ma.ByArtistID)
	if err != nil {
		lg.Errorf("Error GetArtist (%v)", err)
	}

	ma.ByArtist.Members, err = ma.ByArtist.GetMembers(db, *ma.ByArtistID)
	if err != nil {
		lg.Errorf("Error Scanning Persons (%v)", err)
	}

	return ma, nil
}

// GetTracks gets all of the MusicRecording objects for a MusicAlbum
func (obj *MusicAlbum) GetTracks(db *sqlx.DB, id IDType) ([]*MusicRecording, error) {
	// NOTE - Position is now coming from the musicalbum_tracks table
	query := &SelectQuery{}
	query.SetQuery(`
    SELECT musicrecording.id, musicrecording.cid, musicrecording.type, musicrecording.context,
        musicrecording.created_at, musicrecording.updated_at, musicrecording.name, musicrecording.duration,
        musicrecording.isrc, musicalbum_tracks.position, musicrecording.genres, musicrecording.audio_id,
        musicrecording.by_artist_id, musicrecording.recording_of_id
    FROM musicrecording
    INNER JOIN musicalbum_tracks ON musicalbum_tracks.musicrecording_id = musicrecording.id
    WHERE musicalbum_tracks.musicalbum_id=$1
    `, id)

	tracks := []*MusicRecording{}
	err := db.Select(&tracks, query.Query(), query.Args()...)
	if err != nil {
		return nil, fmt.Errorf("Error GetTracks (%v)", err)
	}

	for _, track := range tracks {
		track.Audio, err = track.GetAudio(db, *track.AudioID)
		if err != nil {
			lg.Errorf("Error Scanning AudioObject (%v)", err)
		}

		track.Rights, err = track.GetRights(db, track.ID)
		if err != nil {
			lg.Errorf("Error Scanning Rights (%v)", err)
		}

		track.ByArtist, err = (&MusicGroup{}).Get(db, *track.ByArtistID)
		if err != nil {
			lg.Errorf("Error Scanning Rights (%v)", err)
		}

		track.RecordingOf, err = track.GetComposition(db, *track.RecordingOfID)
		if err != nil {
			lg.Errorf("Error Scanning MusicComposition (%v)", err)
		}

		track.RecordingOf.Composer, err = track.RecordingOf.GetComposer(db, *track.RecordingOfID)
		if err != nil {
			lg.Errorf("Error getting composer (%v)", err)
		}
	}

	return tracks, nil
}

// GetArtist the MusicGroup object for a MusicAlbum
func (obj *MusicAlbum) GetArtist(db *sqlx.DB, musicGroupID IDType) (*MusicGroup, error) {
	query := &SelectQuery{}
	query.SetQuery(`
    SELECT *
    FROM musicgroup
    WHERE musicgroup.id=$1
    `, musicGroupID)

	musicgroup := &MusicGroup{}
	err := db.Get(musicgroup, query.Query(), query.Args()...)
	if err != nil {
		return nil, fmt.Errorf("Error getting MusicAlbum artist (%v)", err)
	}

	musicgroup.Image, err = musicgroup.GetImage(db, *musicgroup.ImageID)
	if err != nil {
		lg.Errorf("Error getting MusicGroup ImageObject (%v)", err)
	}

	return musicgroup, nil
}
