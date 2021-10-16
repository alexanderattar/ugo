package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/consensys/ugo/pkg/lg"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// MusicPlaylist model
type MusicPlaylist struct {
	ID        IDType            `json:"id"`
	CID       string            `json:"cid"`
	CIDs      pq.StringArray    `json:"cids"`
	Type      string            `json:"@type"`
	Context   string            `json:"@context"`
	CreatedAt time.Time         `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time         `json:"updatedAt" db:"updated_at"`
	Name      string            `json:"name"`
	Tracks    []*MusicRecording `json:"tracks"`
	Image     *ImageObject      `json:"image"`
	ImageID   *IDType           `json:"-" db:"image_id"`
	ByUser    *Person           `json:"byUser"`
	ByUserID  *IDType           `json:"-" db:"by_user_id"`
}

// All gets all of the given MusicRelease objects
func (obj *MusicPlaylist) All(db *sqlx.DB, byUser *IDType, query *SelectQuery) ([]*MusicPlaylist, error) {
	if query.Limit == 0 {
		query.Limit = 20
	}

	fmt.Println("here", byUser)

	if byUser != nil { // filter the releases by artist cid
		query.SetQuery(`
        SELECT musicplaylist.id, musicplaylist.cid, musicplaylist.type,
			musicplaylist.context, musicplaylist.name, musicplaylist.context,
			musicplaylist.created_at, musicplaylist.updated_at, musicplaylist.by_user_id,
            musicplaylist.image_id
        FROM musicplaylist
		JOIN person ON musicplaylist.by_user_id = person.id
		WHERE musicplaylist.by_user_id = $1
        `, *byUser)

	} else {
		query.SetQuery(`
        SELECT * FROM musicplaylist
        `)
	}

	musicplaylists := []*MusicPlaylist{}
	err := db.Select(&musicplaylists, query.Query(), query.Args()...)
	if err != nil {
		return nil, fmt.Errorf("Error getting all MusicPlaylists (%v)", err)
	}

	// TODO - Optimize query to grab all fields instead of round trips to the db for subobjects
	// Serialize the musicalbum and image
	for _, musicplaylist := range musicplaylists {
		musicplaylist.Tracks, err = (&MusicPlaylist{}).GetTracks(db, musicplaylist.ID)
		if err != nil {
			lg.Errorf("Error getting musicplaylist.Tracks (%v)", err)
		}

		musicplaylist.Image, err = (&ImageObject{}).Get(db, *musicplaylist.ImageID)
		if err != nil {
			lg.Errorf("Error getting musicplaylist.Image (%v)", err)
		}

		musicplaylist.ByUser, err = (&Person{}).Get(db, *musicplaylist.ByUserID, "")
		if err != nil {
			lg.Errorf("Error getting musicplaylist.ByUser (%v)", err)
		}
	}
	return musicplaylists, err
}

// Create a MusicPlaylist
func (obj *MusicPlaylist) Create(db *sqlx.DB) (musicplaylistID IDType, err error) {
	err = obj.validate()
	if err != nil {
		return 0, err
	}

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	var CIDs []string
	// On initial POST from portal, there is no CID
	if obj.CID != "" {
		CIDs = []string{obj.CID}
	}

	var imageID interface{}

	if obj.Image != nil {
		// Create ImageObject
		err = tx.QueryRow(
			`INSERT INTO imageobject (
				cid, type, context, created_at, updated_at,
				content_url, encoding_format
			) VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id`,
			obj.Image.CID, obj.Image.Type, obj.Image.Context, time.Now(), time.Now(),
			obj.Image.ContentURL, obj.Image.EncodingFormat,
		).Scan(&imageID)

		if err != nil {
			return 0, fmt.Errorf("Error creating playlist image (%v)", err)
		}
	}

	// Create Paylist
	err = tx.QueryRow(
		`INSERT INTO musicplaylist(
			cid, cids, type, context, created_at, updated_at,
			name, by_user_id, image_id
		) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`,
		obj.CID, pq.Array(CIDs), obj.Type, obj.Context, time.Now(), time.Now(),
		obj.Name, obj.ByUser.ID, imageID,
	).Scan(&musicplaylistID)

	if err != nil {
		return 0, fmt.Errorf("Error inserting into musicplaylist (%v)", err)
	}

	// Relate tracks to playlist
	for i, t := range obj.Tracks {
		_, err = tx.Exec(
			`INSERT INTO musicplaylist_tracks(
				musicplaylist_id, musicrecording_id, position
			) VALUES($1, $2, $3)`,
			musicplaylistID, t.ID, i+1,
		)

		if err != nil {
			return 0, fmt.Errorf("Error adding track to playlist (%v)", err)
		}
	}

	return musicplaylistID, nil
}

// Update a Playlist
func (obj *MusicPlaylist) Update(db *sqlx.DB, id IDType) (musicplaylistID IDType, err error) {
	err = obj.validate()
	if err != nil {
		return 0, err
	}

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	for _, t := range obj.Tracks {
		// Update playlist track
		_, err = tx.Exec(
			`UPDATE musicplaylist_tracks
			SET position=$1
			WHERE musicplaylist_id=$2 AND musicrecording_id=$3`,
			t.Position, obj.ID, t.ID,
		)
	}

	if obj.Image == nil {
		// Update playlist without ImageObject
		// We can remove this if an ImageObject is required by Playlist
		err = tx.QueryRow(
			`UPDATE musicplaylist
			 SET cid=$1, cids=array_append(cids, $2), type=$3, context=$4, updated_at=$5,
			 name=$6, by_user=$8,
			 WHERE id=$10
			 RETURNING id`,
			obj.CID, obj.CID, obj.Type, obj.Context, time.Now(),
			obj.Name, obj.ByUser.ID,
			id,
		).Scan(&musicplaylistID)

		if err != nil {
			return musicplaylistID, fmt.Errorf("Error updating playlist (%v)", err)
		}
	} else if obj.Image.ID < 1 {
		// Create ImageObject
		_, err = tx.Exec(
			`INSERT INTO imageobject (
				cid, type, context, created_at, updated_at,
				content_url, encoding_format
			) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			obj.Image.CID, obj.Image.Type, obj.Image.Context, time.Now(), time.Now(),
			obj.Image.ContentURL, obj.Image.EncodingFormat,
		)

		if err != nil {
			return musicplaylistID, fmt.Errorf("Error updating image in playlist (%v)", err)
		}

		// Update Playlist with new ImageObject
		err = tx.QueryRow(
			`UPDATE musicplaylist
			 SET cid=$1, cids=array_append(cids, $2), type=$3, context=$4, updated_at=$5,
			 name=$6, by_user_id=$7,
			 image_id=currval('imageobject_id_seq')
			 WHERE id=$8
			 RETURNING id`,
			obj.CID, obj.CID, obj.Type, obj.Context, time.Now(),
			obj.Name, obj.ByUserID,
			id,
		).Scan(&musicplaylistID)

		if err != nil {
			return musicplaylistID, fmt.Errorf("Error updating playlist (%v)", err)
		}
	} else if obj.Image.ID > 0 {
		// Update Image
		_, err = tx.Exec(
			`UPDATE imageobject
			 SET cid=$1, type=$2, context=$3,
			 updated_at=$4, content_url=$5, encoding_format=$6
			 WHERE id=$7`,
			obj.Image.CID, obj.Type, obj.Context, time.Now(),
			obj.Image.ContentURL, obj.Image.EncodingFormat, obj.Image.ID,
		)

		if err != nil {
			return musicplaylistID, fmt.Errorf("Error updating image in playlist (%v)", err)
		}

		// Update Playlist
		err = tx.QueryRow(
			`UPDATE musicplaylist
			 SET cid=$1, cids=array_append(cids, $2), type=$3, context=$4, updated_at=$5,
			 name=$6, by_user_id=$7
			 WHERE id=$8
			 RETURNING id`,
			obj.CID, obj.CID, obj.Type, obj.Context, time.Now(),
			obj.Name, obj.ByUser.ID, id,
		).Scan(&musicplaylistID)

		if err != nil {
			return 0, fmt.Errorf("Error updating playlist (%v)", err)
		}
	}
	return musicplaylistID, nil
}

// Get MusicPlaylist by ID
func (obj *MusicPlaylist) Get(db *sqlx.DB, id IDType) (*MusicPlaylist, error) {
	query := &SelectQuery{}
	query.SetQuery(`
    SELECT *
    FROM musicplaylist
    WHERE musicplaylist.id=$1
    `, id)

	musicplaylist := &MusicPlaylist{}
	err := db.Get(musicplaylist, query.Query(), query.Args()...)
	if err != nil {
		return nil, fmt.Errorf("Error getting MusicPlaylist (%v)", err)
	}

	imageobject := &ImageObject{}
	musicplaylist.Image, err = imageobject.Get(db, *musicplaylist.ImageID)
	if err != nil {
		lg.Errorf("Error getting musicplaylist.Image (%v)", err)
	}

	person := &Person{}
	musicplaylist.ByUser, err = person.Get(db, *musicplaylist.ByUserID, "")
	if err != nil {
		lg.Errorf("Error getting musicplaylist.ByUser (%v)", err)
	}

	musicplaylist.Tracks, err = musicplaylist.GetTracks(db, musicplaylist.ID)
	if err != nil {
		lg.Errorf("Error getting musicplaylist.Tracks (%v)", err)
	}

	return musicplaylist, err
}

// Delete a MusicPlaylist
func (obj *MusicPlaylist) Delete(db *sqlx.DB, musicplaylistID IDType) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	_, err = db.Exec(
		`DELETE FROM musicplaylist
         WHERE musicplaylist.id = $1`,
		musicplaylistID,
	)

	if err != nil {
		return fmt.Errorf("Error deleting playlist (%v)", err)
	}
	return nil
}

// AddTrack adds a Track to a MusicPlaylist
func (obj *MusicPlaylist) AddTrack(db *sqlx.DB, args ...interface{}) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() error {
		if err != nil {
			tx.Rollback()
			return err
		}
		err = tx.Commit()
		if err != nil {
			return err
		}
		return nil
	}()

	// Make sure track exists
	query := &SelectQuery{}
	query.SetQuery(`
	SELECT *
	FROM musicrecording
	WHERE musicrecording.id=$1
	`, args[1])

	musicrecording := &MusicRecording{}
	err = db.Get(musicrecording, query.Query(), query.Args()...)

	if err != nil {
		return fmt.Errorf("Error finding track (%v)", err)
	}

	// Get number of tracks on playlist
	query.SetQuery(`
    SELECT count(*)
	FROM musicplaylist_tracks
	WHERE musicplaylist_id=$1
    `, args[0])

	var count int
	err = db.Get(&count, query.Query(), query.Args()...)
	if err != nil {
		return fmt.Errorf("Error getting playlist track count (%v)", err)
	}

	if err != nil {
		return fmt.Errorf("Error in AddTrack (%v)", err)
	}

	position := count + 1
	// Relate track to playlist
	_, err = tx.Exec(
		`INSERT INTO musicplaylist_tracks(
			musicplaylist_id, musicrecording_id, position
		) VALUES($1, $2, $3)`,
		args[0], args[1], position,
	)

	if err != nil {
		return fmt.Errorf("Error in AddTrack (%v)", err)
	}

	return err
}

// DeleteTrack deletes a Track from a MusicPlaylist
func (obj *MusicPlaylist) DeleteTrack(db *sqlx.DB, args ...interface{}) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() error {
		if err != nil {
			tx.Rollback()
			return err
		}
		err = tx.Commit()
		if err != nil {
			return err
		}
		return nil
	}()

	// Get position for track position updates
	query := &SelectQuery{}
	query.SetQuery(`
	SELECT position
	FROM musicplaylist_tracks
	WHERE musicplaylist_id=$1 AND musicrecording_id=$2
	`, args[0], args[1])

	var position int
	err = db.Get(&position, query.Query(), query.Args()...)
	if err != nil {
		return fmt.Errorf("Error getting playlist track count (%v)", err)
	}

	// Delete track
	result, err := tx.Exec(`
		DELETE FROM musicplaylist_tracks
		WHERE musicplaylist_id=$1 AND musicrecording_id=$2`,
		args[0], args[1],
	)

	rows, err := result.RowsAffected()
	if rows == 0 {
		return errors.New("No entry found")
	}

	if err != nil {
		return fmt.Errorf("Error in DeleteTrack (%v)", err)
	}

	// Decrement position fields on tracks with larger position values than the track deleted
	tracks, err := obj.GetTracks(db, args[0].(int64))
	if err != nil {
		return fmt.Errorf("Error getting tracks (%v)", err)
	}

	for _, t := range tracks {
		if *t.Position > position {
			_, err = tx.Exec(
				`UPDATE musicplaylist_tracks SET position=$1 WHERE musicplaylist_id=$2 AND musicrecording_id=$3`,
				*t.Position-1, args[0], t.ID,
			)
		}

		if err != nil {
			return fmt.Errorf("Error updating track (%v)", err)
		}
	}

	return err
}

// GetTracks gets all of the MusicRecording objects for a MusicPlaylist
func (obj *MusicPlaylist) GetTracks(db *sqlx.DB, id IDType) ([]*MusicRecording, error) {
	query := &SelectQuery{}
	query.OrderBy = "musicplaylist_tracks.position"
	query.SetQuery(`
    SELECT musicrecording.id, musicrecording.cid, musicrecording.type, musicrecording.context,
        musicrecording.created_at, musicrecording.updated_at, musicrecording.name, musicrecording.duration,
        musicrecording.isrc, musicplaylist_tracks.position, musicrecording.genres, musicrecording.audio_id,
        musicrecording.by_artist_id, musicrecording.recording_of_id
    FROM musicrecording
    INNER JOIN musicplaylist_tracks ON musicplaylist_tracks.musicrecording_id = musicrecording.id
    WHERE musicplaylist_tracks.musicplaylist_id=$1
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

func (obj MusicPlaylist) validate() error {
	return validation.Errors{
		"CID":     validation.Validate(&obj.CID, validation.Required),
		"Context": validation.Validate(&obj.Context, validation.Required),
		"Type":    validation.Validate(&obj.Type, validation.Required),
		"Name":    validation.Validate(&obj.Name, validation.Required),
	}.Filter()
}
