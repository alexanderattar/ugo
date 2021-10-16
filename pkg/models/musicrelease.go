package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/consensys/ugo/pkg/lg"
	"github.com/consensys/ugo/pkg/utils"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// MusicRelease model
type MusicRelease struct {
	ID                 IDType        `json:"id"`
	CID                string        `json:"cid"`
	Type               string        `json:"@type"`
	Context            string        `json:"@context"`
	CreatedAt          time.Time     `json:"createdAt" db:"created_at"`
	UpdatedAt          time.Time     `json:"updatedAt" db:"updated_at"`
	Description        *string       `json:"description"`
	DatePublished      *string       `json:"datePublished" db:"date_published"`
	CatalogNumber      *string       `json:"catalogNumber" db:"catalog_number"`
	MusicReleaseFormat *string       `json:"musicReleaseFormat" db:"music_release_format"`
	Price              *float64      `json:"price"`
	RecordLabel        *Organization `json:"recordLabel" db:"record_label"`
	ReleaseOf          *MusicAlbum   `json:"releaseOf" db:"release_of"`
	Image              *ImageObject  `json:"image"`
	RecordingLabelID   *IDType       `json:"-" db:"record_label_id"`
	ReleaseOfID        *IDType       `json:"-" db:"release_of_id"`
	ImageID            *IDType       `json:"-" db:"image_id"`
	Visibility         *string       `json:"-"`
	Active             bool          `json:"-"`
}

// All gets all of the given MusicRelease objects
func (obj *MusicRelease) All(db *sqlx.DB, byArtist *IDType, query *SelectQuery) ([]*MusicRelease, error) {
	if query.Limit == 0 {
		query.Limit = 20
	}

	if byArtist != nil { // filter the releases by artist cid
		query.SetQuery(`
        SELECT musicrelease.id, musicrelease.cid, musicrelease.type,
            musicrelease.context, musicrelease.created_at, musicrelease.updated_at,
            musicrelease.active, musicrelease.description, musicrelease.date_published,
            musicrelease.catalog_number, musicrelease.music_release_format, musicrelease.price,
            musicrelease.record_label_id, musicrelease.release_of_id, musicrelease.image_id
        FROM musicrelease
        JOIN musicalbum ON musicrelease.release_of_id = musicalbum.id
        JOIN musicgroup ON musicalbum.by_artist_id = musicgroup.id
        WHERE musicgroup.id = $1 AND musicrelease.active = true
        `, *byArtist)

		// Get db field that corresponds to the value passed
		query.OrderBy = utils.ParseMusicReleaseOrderBy(query.OrderBy)

	} else if query.OrderBy != "" {
		query.SetQuery(`
        SELECT musicrelease.id, musicrelease.cid, musicrelease.type,
            musicrelease.context, musicrelease.created_at, musicrelease.updated_at,
            musicrelease.active, musicrelease.description, musicrelease.date_published,
            musicrelease.catalog_number, musicrelease.music_release_format, musicrelease.price,
            musicrelease.record_label_id, musicrelease.release_of_id, musicrelease.image_id
        FROM musicrelease
        JOIN musicalbum ON musicrelease.release_of_id = musicalbum.id
        JOIN musicgroup ON musicalbum.by_artist_id = musicgroup.id
        WHERE musicrelease.active = true
        `)

		// Get db field that corresponds to the value passed
		query.OrderBy = utils.ParseMusicReleaseOrderBy(query.OrderBy)

	} else {
		query.SetQuery(`
        SELECT *
        FROM musicrelease
        WHERE musicrelease.active = true
        `)
	}

	musicreleases := []*MusicRelease{}
	err := db.Select(&musicreleases, query.Query(), query.Args()...)
	if err != nil {
		return nil, fmt.Errorf("Error getting all MusicReleases (%v)", err)
	}

	// TODO - Optimize query to grab all fields instead of round trips to the db for subobjects
	// Serialize the musicalbum and image
	for _, musicrelease := range musicreleases {
		musicrelease.ReleaseOf, err = (&MusicAlbum{}).Get(db, *musicrelease.ReleaseOfID)
		if err != nil {
			lg.Errorf("Error Scanning MusicAlbum (%v)", err)
		}

		musicrelease.Image, err = musicrelease.GetImage(db, *musicrelease.ImageID)
		if err != nil {
			lg.Errorf("Error Scanning ImageObject (%v)", err)
		}
	}
	return musicreleases, err
}

// AllInactive gets all of the inactive MusicRelease objects
func (obj *MusicRelease) AllInactive(db *sqlx.DB, query *SelectQuery) ([]*MusicRelease, error) {
	if query.Limit == 0 {
		query.Limit = 50
	}

	query.SetQuery(`
        SELECT *
        FROM musicrelease
        WHERE musicrelease.active = false
    `)

	musicreleases := []*MusicRelease{}
	err := db.Select(&musicreleases, query.Query(), query.Args()...)
	if err != nil {
		return nil, fmt.Errorf("Error getting all MusicReleases (%v)", err)
	}

	// TODO - Optimize query to grab all fields instead of round trips to the db for subobjects
	// Serialize the musicalbum and image
	for _, musicrelease := range musicreleases {
		musicrelease.ReleaseOf, err = (&MusicAlbum{}).Get(db, *musicrelease.ReleaseOfID)
		if err != nil {
			lg.Errorf("Error Scanning MusicAlbum (%v)", err)
		}

		musicrelease.Image, err = musicrelease.GetImage(db, *musicrelease.ImageID)
		if err != nil {
			lg.Errorf("Error Scanning ImageObject (%v)", err)
		}
	}
	return musicreleases, err
}

// Get MusicRelease by ID
func (obj *MusicRelease) GetByID(db *sqlx.DB, id IDType) (*MusicRelease, error) {
	query := &SelectQuery{}
	query.SetQuery(`
    SELECT *
    FROM musicrelease
    WHERE musicrelease.id=$1
    `, id)

	musicrelease := &MusicRelease{}
	err := db.Get(musicrelease, query.Query(), query.Args()...)
	if err != nil {
		return nil, fmt.Errorf("Error getting MusicRelease (%v)", err)
	}

	musicrelease.ReleaseOf, err = (&MusicAlbum{}).Get(db, *musicrelease.ReleaseOfID)
	if err != nil {
		lg.Errorf("Error Scanning MusicAlbum (%v)", err)
	}

	musicrelease.Image, err = musicrelease.GetImage(db, *musicrelease.ImageID)
	if err != nil {
		lg.Errorf("Error Scanning ImageObject (%v)", err)
	}

	return musicrelease, err
}

// Create a MusicRelease
func (obj *MusicRelease) Create(db *sqlx.DB) (created *MusicRelease, err error) {
	err = obj.validate()
	if err != nil {
		return nil, err
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Create Album
	_, err = tx.Exec(
		`INSERT INTO musicalbum(
			cid, type, context, created_at, updated_at,
			name, album_production_type, album_release_type, by_artist_id
        ) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		obj.ReleaseOf.CID, obj.ReleaseOf.Type, obj.Context, time.Now(), time.Now(),
		obj.ReleaseOf.Name, obj.ReleaseOf.AlbumProductionType,
		obj.ReleaseOf.AlbumReleaseType, obj.ReleaseOf.ByArtist.ID,
	)

	if err != nil {
		return nil, fmt.Errorf("Error inserting into musicalbum (%v)", err)
	}

	// Create tracks
	for _, t := range obj.ReleaseOf.Tracks {
		err = validateTrack(t)
		if err != nil {
			return nil, err
		}

		// Create Composition
		_, err = tx.Exec(
			`INSERT INTO musiccomposition(
				cid, type, context, created_at, updated_at, name, iswc
			) VALUES($1, $2, $3, $4, $5, $6, $7)`,
			t.RecordingOf.CID, t.RecordingOf.Type, t.Context, time.Now(), time.Now(),
			t.RecordingOf.Name, t.RecordingOf.Iswc,
		)

		if err != nil {
			return nil, fmt.Errorf("Error inserting into musiccomposition (%v)", err)
		}

		for _, composer := range t.RecordingOf.Composer {
			_, err = tx.Exec(`
				INSERT INTO musiccomposition_composer(
					musiccomposition_id, person_id
				) VALUES(currval('musiccomposition_id_seq'), $1)`,
				composer.ID,
			)

			if err != nil {
				return nil, fmt.Errorf("Error inserting into musiccomposition_composer (%v)", err)
			}
		}

		// Create AudioObject
		_, err = tx.Exec(
			`INSERT INTO audioobject(
                cid, type, context, created_at, updated_at, content_url, encoding_format
            ) VALUES($1, $2, $3, $4, $5, $6, $7)`,
			t.Audio.CID, t.Audio.Type, t.Context, time.Now(), time.Now(),
			t.Audio.ContentURL, t.Audio.EncodingFormat,
		)

		if err != nil {
			return nil, fmt.Errorf("Error inserting into audioobject (%v)", err)
		}

		// Create Recording
		_, err = tx.Exec(
			`INSERT INTO musicrecording(
                cid, type, context, created_at, updated_at, name, duration, isrc, position, genres,
                by_artist_id, audio_id, recording_of_id
            ) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11,
            currval('audioobject_id_seq'), currval('musiccomposition_id_seq'))`,
			t.CID, t.Type, t.Context, time.Now(), time.Now(),
			t.Name, t.Duration, t.Isrc, t.Position, pq.Array(t.Genres), obj.ReleaseOf.ByArtist.ID,
		)

		if err != nil {
			return nil, fmt.Errorf("Error inserting into musicrecording (%v)", err)
		}

		// Relate tracks to the MusicAlbum
		_, err = tx.Exec(
			`INSERT INTO musicalbum_tracks(
				musicalbum_id, musicrecording_id
			) VALUES(currval('musicalbum_id_seq'), currval('musicrecording_id_seq'))`,
		)

		if err != nil {
			return nil, fmt.Errorf("Error inserting into musicalbum_tracks (%v)", err)
		}
	}

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
		return nil, fmt.Errorf("Error inserting into imageobject (%v)", err)
	}

	// Create Release
	// TODO - Add currval('organization_id_seq')
	var musicreleaseID IDType

	err = tx.QueryRow(
		`INSERT INTO musicrelease(
			cid, type, context, created_at, updated_at, active, description, date_published,
			catalog_number, music_release_format, price, record_label_id, release_of_id, image_id
		) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, currval('musicalbum_id_seq'), currval('imageobject_id_seq'))
		RETURNING id`,
		obj.CID, obj.Type, obj.Context, time.Now(), time.Now(), true,
		obj.Description, obj.DatePublished, obj.CatalogNumber,
		obj.MusicReleaseFormat, obj.Price, nil,
	).Scan(&musicreleaseID)

	if err != nil {
		return nil, fmt.Errorf("Error inserting into musicrelease (%v)", err)
	}

	obj.ID = musicreleaseID

	return obj, nil
}

// Link associates tracks with a MusicRelease
func (obj *MusicRelease) Link(db *sqlx.DB) (created *MusicRelease, err error) {
	err = obj.validate()
	if err != nil {
		return nil, err
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Create Album
	_, err = tx.Exec(
		`INSERT INTO musicalbum(
			cid, type, context, created_at, updated_at,
			name, album_production_type, album_release_type, by_artist_id
	 	) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`,
		obj.ReleaseOf.CID, obj.ReleaseOf.Type, obj.Context, time.Now(), time.Now(),
		obj.ReleaseOf.Name, obj.ReleaseOf.AlbumProductionType,
		obj.ReleaseOf.AlbumReleaseType, obj.ReleaseOf.ByArtist.ID,
	)

	if err != nil {
		return nil, fmt.Errorf("Error inserting into musicalbum (%v)", err)
	}

	// Create tracks
	for _, t := range obj.ReleaseOf.Tracks {
		err = validateTrack(t)
		if err != nil {
			return nil, err
		}

		// Relate tracks to the MusicAlbum
		_, err = tx.Exec(
			`INSERT INTO musicalbum_tracks(
				musicalbum_id, musicrecording_id, position
			) VALUES(currval('musicalbum_id_seq'), $1, $2)`,
			t.ID, t.Position,
		)

		if err != nil {
			return nil, fmt.Errorf("Error inserting into musicalbum_tracks (%v)", err)
		}
	}

	_, err = tx.Exec(
		`INSERT INTO imageobject (
            cid, type, context, created_at, updated_at,
            content_url, encoding_format
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		obj.Image.CID, obj.Image.Type, obj.Image.Context, time.Now(), time.Now(),
		obj.Image.ContentURL, obj.Image.EncodingFormat,
	)

	if err != nil {
		return nil, fmt.Errorf("Error inserting into imageobject (%v)", err)
	}

	// Create Release
	var musicreleaseID int64

	err = tx.QueryRow(
		`INSERT INTO musicrelease(
			cid, type, context, created_at, updated_at, active, description, date_published,
			catalog_number, music_release_format, price, record_label_id, release_of_id, image_id
		) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, currval('musicalbum_id_seq'), currval('imageobject_id_seq'))
		RETURNING id`,
		obj.CID, obj.Type, obj.Context, time.Now(), time.Now(), true,
		obj.Description, obj.DatePublished, obj.CatalogNumber,
		obj.MusicReleaseFormat, obj.Price, nil,
	).Scan(&musicreleaseID)

	if err != nil {
		return nil, fmt.Errorf("Error inserting into musicrelease (%v)", err)
	}

	obj.ID = musicreleaseID

	return obj, nil
}

// Update a MusicRelease
// Note that PUT requests require id fields to be sent for sub-objects
func (obj *MusicRelease) Update(db *sqlx.DB, id IDType) (updated *MusicRelease, err error) {
	err = obj.validate()
	if err != nil {
		return nil, err
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Update Album
	_, err = tx.Exec(
		`UPDATE musicalbum
		SET cid=$1, type=$2, context=$3, updated_at=$4, name=$5,
		album_production_type=$6, album_release_type=$7, by_artist_id=$8
		WHERE id=$9`,
		obj.ReleaseOf.CID, obj.ReleaseOf.Type, obj.Context, time.Now(),
		obj.ReleaseOf.Name, obj.ReleaseOf.AlbumProductionType,
		obj.ReleaseOf.AlbumReleaseType, obj.ReleaseOf.ByArtist.ID, obj.ReleaseOf.ID,
	)

	if err != nil {
		return nil, fmt.Errorf("Error updating musicalbum (%v)", err)
	}

	// Update tracks
	for _, t := range obj.ReleaseOf.Tracks {
		err = validateTrack(t)
		if err != nil {
			return nil, err
		}

		// Update Composition
		_, err = tx.Exec(
			`UPDATE musiccomposition
			SET cid=$1, type=$2, context=$3, updated_at=$4, name=$5, iswc=$6
			WHERE id=$7`,
			t.CID, t.RecordingOf.Type, t.Context, time.Now(),
			t.RecordingOf.Name, t.RecordingOf.Iswc, t.RecordingOf.ID,
		)

		if err != nil {
			return nil, fmt.Errorf("Error updating musiccomposition (%v)", err)
		}

		// Update AudioObject
		_, err = tx.Exec(
			`UPDATE audioobject
			SET cid=$1, type=$2, context=$3, updated_at=$4, content_url=$5, encoding_format=$6
			WHERE id=$7`,
			t.CID, t.Audio.Type, t.Audio.Context, time.Now(),
			t.Audio.ContentURL, t.Audio.EncodingFormat, t.Audio.ID,
		)

		if err != nil {
			return nil, fmt.Errorf("Error updating audioobject (%v)", err)
		}

		// Update Recording
		_, err = tx.Exec(
			`UPDATE musicrecording
			SET cid=$1, type=$2, context=$3, updated_at=$4, name=$5,
			duration=$6, isrc=$7, position=$8, genres=$9, by_artist_id=$10, audio_id=$11, recording_of_id=$12
			WHERE id=$13`,
			t.CID, t.Type, t.Context, time.Now(),
			t.Name, t.Duration, t.Isrc, t.Position, pq.Array(t.Genres), obj.ReleaseOf.ByArtist.ID,
			t.Audio.ID, t.RecordingOf.ID, t.ID,
		)

		if err != nil {
			return nil, fmt.Errorf("Error updating musicrecording (%v)", err)
		}
	}

	// Update Release
	// TODO - Add currval('organization_id_seq')
	_, err = tx.Exec(
		`UPDATE musicrelease
		SET cid=$1, type=$2, context=$3, updated_at=$4, active=$5, description=$6, date_published=$7,
		catalog_number=$8, music_release_format=$9, price=$10, record_label_id=$11, release_of_id=$12
		WHERE id=$13`,
		obj.CID, obj.Type, obj.Context, time.Now(), obj.Active, obj.Description, obj.DatePublished,
		obj.CatalogNumber, obj.MusicReleaseFormat, obj.Price, nil, obj.ReleaseOf.ID, id,
	)

	if err != nil {
		return nil, fmt.Errorf("Error updating musicrelease (%v)", err)
	}

	// Update ImageObject
	_, err = tx.Exec(
		`UPDATE imageobject
		SET cid=$1, type=$2, context=$3, updated_at=$4,
		content_url=$5, encoding_format=$6
		WHERE id=$7`,
		obj.Image.CID, obj.Image.Type, obj.Image.Context, time.Now(),
		obj.Image.ContentURL, obj.Image.EncodingFormat, obj.Image.ID,
	)

	if err != nil {
		return nil, fmt.Errorf("Error updating imageobject (%v)", err)
	}

	return obj, nil
}

// Delete a MusicRelease
func (obj *MusicRelease) Delete(db *sqlx.DB, id IDType) (err error) {
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

	_, err = db.Exec(`
        DELETE FROM audioobject
            WHERE audioobject.id in (
            SELECT audioobject.id from audioobject
            INNER JOIN musicrecording ON musicrecording.audio_id = audioobject.id
            INNER JOIN musicalbum_tracks ON musicalbum_tracks.musicrecording_id = musicrecording.id
            INNER JOIN musicalbum ON musicalbum.id = musicalbum_tracks.musicalbum_id
            WHERE musicalbum.id=$1
        )`, id,
	)

	if err != nil {
		return fmt.Errorf("Error deleting audioobject (%v)", err)
	}

	// Get this query to work
	_, err = db.Exec(`
		DELETE from musiccomposition_composer 
			WHERE musiccomposition_composer.id in (
			SELECT musiccomposition_composer.id from musiccomposition_composer
			INNER JOIN musiccomposition ON musiccomposition.id = musiccomposition_composer.musiccomposition_id
			INNER JOIN musicrecording ON musicrecording.recording_of_id = musiccomposition.id
			INNER JOIN musicalbum_tracks ON musicalbum_tracks.musicrecording_id = musicrecording.id
			INNER JOIN musicalbum ON musicalbum.id = musicalbum_tracks.musicalbum_id
			WHERE musicalbum.id=$1
		)`, id,
	)

	if err != nil {
		return fmt.Errorf("Error deleting musiccomposition_composer (%v)", err)
	}

	// TODO - Compositions are currently not getting deleted because the joins no longer
	// work after recordings are deleted from the cascade from deleting audio objects
	_, err = db.Exec(`
        DELETE FROM musiccomposition
            WHERE musiccomposition.id in (
            SELECT musiccomposition.id from musiccomposition
            INNER JOIN musicrecording ON musicrecording.recording_of_id = musiccomposition.id
            INNER JOIN musicalbum_tracks ON musicalbum_tracks.musicrecording_id = musicrecording.id
            INNER JOIN musicalbum ON musicalbum.id = musicalbum_tracks.musicalbum_id
            WHERE musicalbum.id=$1
        )`, id,
	)

	if err != nil {
		return fmt.Errorf("Error deleting musiccomposition (%v)", err)
	}

	_, err = db.Exec(
		`DELETE FROM musicrecording
         WHERE musicrecording.id in (
             SELECT musicrecording.id from musicrecording
             INNER JOIN musicalbum_tracks ON musicalbum_tracks.musicrecording_id = musicrecording.id
             INNER JOIN musicalbum ON musicalbum.id = musicalbum_tracks.musicalbum_id
             WHERE musicalbum.id=$1
        )`, id,
	)

	if err != nil {
		return fmt.Errorf("Error updating musicrecording (%v)", err)
	}

	_, err = db.Exec(`
        DELETE from musicalbum where musicalbum.id in (
        SELECT musicalbum.id from musicalbum
        JOIN musicrelease on musicrelease.release_of_id = musicalbum.id
        WHERE musicrelease.release_of_id=$1
        )`, id,
	)

	if err != nil {
		return fmt.Errorf("Error updating musicalbum (%v)", err)
	}

	return nil
}

// GetImage gets the ImageObject for a MusicRelease
func (obj *MusicRelease) GetImage(db *sqlx.DB, imageobjectID IDType) (*ImageObject, error) {
	query := &SelectQuery{}
	query.SetQuery(`
        SELECT imageobject.id, imageobject.cid, imageobject.type, imageobject.context,
            imageobject.created_at, imageobject.updated_at, imageobject.content_url, imageobject.encoding_format
        FROM imageobject
        WHERE id=$1
    `, imageobjectID)

	imageobject := &ImageObject{}
	err := db.Get(imageobject, query.Query(), query.Args()...)
	if err != nil {
		return nil, fmt.Errorf("Error GetImage (%v)", err)
	}
	return imageobject, nil
}

// AddTrack adds a Track to a MusicRelease
func (obj *MusicRelease) AddTrack(db *sqlx.DB, args ...interface{}) error {
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
	query := `
	SELECT *
	FROM musicrecording
	WHERE musicrecording.id=$1
	`

	musicrecording := &MusicRecording{}
	err = db.Get(musicrecording, query, args[1])

	if err != nil {
		return fmt.Errorf("Error finding track (%v)", err)
	}

	// Relate track to album
	_, err = tx.Exec(
		`INSERT INTO musicalbum_tracks(
			musicalbum_id, musicrecording_id
		) VALUES($1, $2)`,
		args[0], args[1],
	)

	if err != nil {
		return fmt.Errorf("Error in AddTrack (%v)", err)
	}

	return err
}

// DeleteTrack deletes a Track from a MusicRelease
func (obj *MusicRelease) DeleteTrack(db *sqlx.DB, args ...interface{}) error {
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

	result, err := db.Exec(`
		DELETE FROM musicalbum_tracks
		WHERE musicalbum_id=$1 AND musicrecording_id=$2`,
		args[0], args[1],
	)

	rows, err := result.RowsAffected()
	if rows == 0 {
		return errors.New("No entry found")
	}

	if err != nil {
		return fmt.Errorf("Error in DeleteTrack (%v)", err)
	}

	return err
}

func (obj MusicRelease) validate() error {
	return validation.Errors{
		"CID":               validation.Validate(&obj.CID, validation.Required),
		"Context":           validation.Validate(&obj.Context, validation.Required),
		"Type":              validation.Validate(&obj.Type, validation.Required),
		"ReleaseOf CID":     validation.Validate(&obj.ReleaseOf.CID, validation.Required),
		"ReleaseOf Type":    validation.Validate(&obj.ReleaseOf.Type, validation.Required),
		"ReleaseOf Context": validation.Validate(&obj.ReleaseOf.Context, validation.Required),
	}.Filter()
}

func validateTrack(track *MusicRecording) error {
	return validation.Errors{
		"Track CID":              validation.Validate(track.CID, validation.Required),
		"Track Context":          validation.Validate(track.Context, validation.Required),
		"Track Type":             validation.Validate(track.Type, validation.Required),
		"Track Audio ContentURL": validation.Validate(track.Audio.ContentURL, validation.Required),
	}.Filter()
}
