package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/consensys/ugo/pkg/lg"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// MusicRecording model
type MusicRecording struct {
	ID            IDType            `json:"id"`
	CID           string            `json:"cid"`
	Type          string            `json:"@type"`
	Context       string            `json:"@context"`
	CreatedAt     time.Time         `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time         `json:"updatedAt" db:"updated_at"`
	Name          *string           `json:"name"`
	RecordingOf   *MusicComposition `json:"recordingOf"`
	ByArtist      *MusicGroup       `json:"byArtist"`
	Duration      *string           `json:"duration"`
	Isrc          *string           `json:"isrc"`
	Position      *int              `json:"position"`
	Genres        pq.StringArray    `json:"genres"`
	Audio         *AudioObject      `json:"audio"`
	Image         *ImageObject      `json:"image"`
	Rights        []*Right          `json:"rights"`
	AudioID       *IDType           `json:"-" db:"audio_id"`
	ByArtistID    *IDType           `json:"-" db:"by_artist_id"`
	RecordingOfID *IDType           `json:"-" db:"recording_of_id"`
	ImageID       *IDType           `json:"-" db:"image_id"`
	Visibility    *string           `json:"-"`
}

// All gets all of the given MusicGroup objects
func (obj *MusicRecording) All(db *sqlx.DB, byArtistID *IDType, query *SelectQuery) ([]*MusicRecording, error) {
	if byArtistID != nil {
		query.SetQuery(`
        SELECT musicrecording.id, musicrecording.cid, musicrecording.type,
			musicrecording.context, musicrecording.created_at, musicrecording.updated_at,
			musicrecording.name, musicrecording.duration, musicrecording.isrc,
			musicrecording.position, musicrecording.genres, musicrecording.audio_id,
			musicrecording.image_id, musicrecording.recording_of_id, musicrecording.by_artist_id
        FROM musicrecording
        INNER JOIN musiccomposition ON musiccomposition.id = musicrecording.recording_of_id
        INNER JOIN musicgroup ON musicgroup.id = musicrecording.by_artist_id
        INNER JOIN audioobject ON audioobject.id = musicrecording.audio_id
		LEFT JOIN imageobject ON imageobject.id = musicrecording.image_id
		WHERE musicrecording.by_artist_id = $1
		`, *byArtistID)
	} else {
		query.SetQuery(`
        SELECT *
        FROM musicrecording
        `)
	}

	musicrecordings := []*MusicRecording{}
	err := db.Select(&musicrecordings, query.Query(), query.Args()...)
	if err != nil {
		return nil, fmt.Errorf("Error getting all MusicRecordings (%v)", err)
	}

	for _, musicrecording := range musicrecordings {
		audio, err := musicrecording.GetAudio(db, *musicrecording.AudioID)

		if err != nil {
			lg.Errorf("Error Scanning AudioObject (%v)", err)
		}

		musicrecording.Audio = audio
		composition, err := musicrecording.GetComposition(db, *musicrecording.RecordingOfID)

		if err != nil {
			lg.Errorf("Error Scanning MusicComposition (%v)", err)
		}

		composition.Composer, err = composition.GetComposer(db, composition.ID)
		musicrecording.RecordingOf = composition

		musicgroup := &MusicGroup{}
		musicgroup, err = musicgroup.Get(db, *musicrecording.ByArtistID)

		if err != nil {
			lg.Errorf("Error Scanning MusicGroup (%v)", err)
		}

		musicrecording.ByArtist = musicgroup

		if musicrecording.ImageID != nil && *musicrecording.ImageID > 0 {
			imageObject := &ImageObject{}
			imageObject, err = imageObject.Get(db, *musicrecording.ImageID)

			if err != nil {
				return nil, fmt.Errorf("Error getting MusicRecording Image (%v)", err)
			}

			musicrecording.Image = imageObject
		}

		rights, err := musicrecording.GetRights(db, musicrecording.ID)
		if err != nil {
			lg.Errorf("Error Scanning Rights (%v)", err)
		}

		musicrecording.Rights = rights
	}

	return musicrecordings, nil
}

// Get MusicRecording by ID
func (obj *MusicRecording) Get(db *sqlx.DB, id IDType) (*MusicRecording, error) {
	query := &SelectQuery{}
	query.SetQuery(`
    SELECT *
    FROM musicrecording
    WHERE musicrecording.id=$1
    `, id)

	musicrecording := &MusicRecording{}
	err := db.Get(musicrecording, query.Query(), query.Args()...)

	composition, err := musicrecording.GetComposition(db, *musicrecording.RecordingOfID)
	if err != nil {
		lg.Errorf("Error Scanning MusicComposition (%v)", err)
	}

	composition.Composer, err = composition.GetComposer(db, composition.ID)
	musicrecording.RecordingOf = composition

	if musicrecording.ImageID != nil && *musicrecording.ImageID > 0 {
		imageObject := &ImageObject{}
		imageObject, err = imageObject.Get(db, *musicrecording.ImageID)

		if err != nil {
			return nil, fmt.Errorf("Error getting MusicRecording Image (%v)", err)
		}

		musicrecording.Image = imageObject
	}

	audio, err := musicrecording.GetAudio(db, *musicrecording.AudioID)
	if err != nil {
		lg.Errorf("Error Scanning AudioObject (%v)", err)
	}

	musicrecording.Audio = audio

	// attach the audio object onto the music recording
	audio, err = musicrecording.GetAudio(db, *musicrecording.AudioID)
	if err != nil {
		lg.Errorf("Error Scanning AudioObject (%v)", err)
	}
	musicrecording.Audio = audio

	// attach the byArtist object onto the music recording
	musicgroup := &MusicGroup{}
	musicgroup, err = musicgroup.Get(db, *musicrecording.ByArtistID)

	if err != nil {
		lg.Errorf("Error Scanning MusicGroup (%v)", err)
	}

	musicrecording.ByArtist = musicgroup

	rights, err := musicrecording.GetRights(db, musicrecording.ID)
	if err != nil {
		lg.Errorf("Error Scanning Rights (%v)", err)
	}

	musicrecording.Rights = rights

	return musicrecording, nil
}

// Create a MusicRecording
func (obj *MusicRecording) Create(db *sqlx.DB) (created *MusicRecording, err error) {
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

	var recordingOfID IDType
	err = tx.QueryRow(`
		INSERT INTO musiccomposition (
            cid, type, context, created_at, updated_at, name
        ) VALUES($1, $2, $3, $4, $5, $6)
        RETURNING id`,
		obj.RecordingOf.CID, obj.RecordingOf.Type, obj.RecordingOf.Context, time.Now(), time.Now(),
		obj.RecordingOf.Name,
	).Scan(&recordingOfID)

	if err != nil {
		return nil, fmt.Errorf("Error inserting into musiccomposition (%v)", err)
	}

	for _, composer := range obj.RecordingOf.Composer {
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

	var audioID IDType
	err = tx.QueryRow(`
        INSERT INTO audioobject(
            cid, type, context, created_at, updated_at,
            content_url, encoding_format
        ) VALUES($1, $2, $3, $4, $5, $6, $7)
        RETURNING id`,
		obj.Audio.CID, obj.Audio.Type, obj.Context, time.Now(), time.Now(),
		obj.Audio.ContentURL, obj.Audio.EncodingFormat,
	).Scan(&audioID)

	if err != nil {
		return nil, fmt.Errorf("Error inserting into audioobject (%v)", err)
	}

	var imageID *IDType
	if obj.Image != nil {
		err = tx.QueryRow(`
			INSERT INTO imageobject (
                cid, type, context, created_at, updated_at,
                content_url, encoding_format
            ) VALUES ($1, $2, $3, $4, $5, $6, $7)
            RETURNING id`,
			obj.Image.CID, obj.Image.Type, obj.Context, time.Now(),
			time.Now(), obj.Image.ContentURL, obj.Image.EncodingFormat,
		).Scan(&imageID)

		if err != nil {
			return nil, fmt.Errorf("Error inserting into imageobject (%v)", err)
		}
	}

	var musicrecordingID IDType

	// This is a hack to ensure both created and updated time are exactly the same
	// in order to support the frontend which currently compares these values to display/hide
	// warning icons. We should remove this in favor of the frontend handling the FE logic more appropriately
	createdAndUpdatedTime := time.Now()
	err = tx.QueryRow(`
        INSERT INTO musicrecording(
            cid, type, context, created_at, updated_at, name, duration, isrc, position, genres,
            by_artist_id, audio_id, recording_of_id, image_id
        ) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, currval('audioobject_id_seq'),
        currval('musiccomposition_id_seq'), $12)
        RETURNING id`,
		obj.CID, obj.Type, obj.Context, createdAndUpdatedTime, createdAndUpdatedTime, obj.Name, obj.Duration,
		obj.Isrc, obj.Position, pq.Array(obj.Genres), obj.ByArtist.ID, imageID,
	).Scan(&musicrecordingID)

	if err != nil {
		return nil, fmt.Errorf("Error inserting into musicrecording (%v)", err)
	}

	obj.ID = musicrecordingID

	var rightID IDType
	var totalShares float64
	for _, right := range obj.Rights {
		totalShares += right.PercentageShares

		err = tx.QueryRow(`
			INSERT INTO "right"(
				cid, type, context, created_at, updated_at, percentage_shares, valid_from, valid_through,
				musicrecording_id, person_id
			) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING id`,
			right.CID, right.Type, right.Context, time.Now(), time.Now(), right.PercentageShares, right.ValidFrom,
			right.ValidThrough, obj.ID, right.Party.ID,
		).Scan(&rightID)

		if err != nil {
			return nil, fmt.Errorf("Error inserting into right (%v)", err)
		}
	}

	if totalShares != 100 {
		return nil, fmt.Errorf("Error creating musicrecording: total shares != 100%%")
	}
	return obj, nil
}

// Update a MusicRecording
// Note that PUT requests require id fields to be sent for sub-objects
func (obj *MusicRecording) Update(db *sqlx.DB, musicRecordingID IDType) (updated *MusicRecording, err error) {
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

	imageID, err := obj.UpdateImage(db, tx)

	if err != nil {
		return nil, fmt.Errorf("Error updating imageobject (%v)", err)
	}

	_, err = tx.Exec(`
		UPDATE audioobject
		SET cid=$1, type=$2, context=$3,
		updated_at=$4, content_url=$5, encoding_format=$6
		WHERE id=$7`,
		obj.Audio.CID, obj.Audio.Type, obj.Context, time.Now(),
		obj.Audio.ContentURL, obj.Audio.EncodingFormat, obj.Audio.ID,
	)

	if err != nil {
		return nil, fmt.Errorf("Error updating audioobject (%v)", err)
	}

	_, err = tx.Exec(`
		UPDATE musiccomposition
		SET cid=$1,  type=$2, context=$3, updated_at=$4, name=$5
		WHERE id=$6`,
		obj.RecordingOf.CID, obj.RecordingOf.Type, obj.Context, time.Now(),
		obj.RecordingOf.Name, obj.RecordingOf.ID,
	)

	if err != nil {
		return nil, fmt.Errorf("Error updating musiccomposition (%v)", err)
	}

	_, err = tx.Exec(
		`UPDATE musicrecording
		SET cid=$1, updated_at=$2, type=$3,
		context=$4, name=$5, duration=$6, isrc=$7, position=$8, genres=$9, audio_id=$10,
		by_artist_id=$11, recording_of_id=$12, image_id=$13
		WHERE id=$14`,
		obj.CID, time.Now(), obj.Type,
		obj.Context, obj.Name, obj.Duration, obj.Isrc, obj.Position,
		pq.Array(obj.Genres), obj.Audio.ID, obj.ByArtist.ID, obj.RecordingOf.ID,
		imageID, musicRecordingID,
	)

	if err != nil {
		return nil, fmt.Errorf("Error updating musicrecording (%v)", err)
	}

	oldRights, err := obj.GetRights(db, musicRecordingID)

	// Delete right objects that've been removed
	for _, oldRight := range oldRights {
		var match bool
		for _, r := range obj.Rights {
			if r.ID == oldRight.ID {
				match = true
			}
		}

		if match == false {
			_, err = tx.Exec(
				`DELETE FROM "right"
				 WHERE id = $1`,
				oldRight.ID,
			)

			if err != nil {
				return nil, fmt.Errorf("Error deleting right (%v)", err)
			}
		}
	}

	// Create a right if it doesn't exist yet
	// Update it otherwise
	for _, right := range obj.Rights {
		// Right doesn't exist
		if right.ID == 0 {
			var rightID IDType
			err = tx.QueryRow(`
				INSERT INTO "right"(
					cid, type, context, created_at, updated_at, percentage_shares,
					valid_from, valid_through, musicrecording_id, person_id
				) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
				RETURNING id`,
				right.CID, right.Type, right.Context, time.Now(), time.Now(), right.PercentageShares, right.ValidFrom,
				right.ValidThrough, obj.ID, right.Party.ID,
			).Scan(&rightID)

			if err != nil {
				return nil, fmt.Errorf("Error inserting into right (%v)", err)
			}
		} else {
			// Right exists so update it
			_, err = tx.Exec(`
				UPDATE "right"
				SET cid=$1, type=$2, context=$3, updated_at=$4, percentage_shares=$5, valid_from=$6, valid_through=$7,
				musicrecording_id=$8, person_id=$9
				WHERE id=$10`,
				right.CID, right.Type, right.Context, time.Now(), right.PercentageShares, right.ValidFrom,
				right.ValidThrough, obj.ID, right.Party.ID, right.ID,
			)

			if err != nil {
				return nil, fmt.Errorf("Error updating right (%v)", err)
			}
		}
	}

	return obj, nil
}

// UpdateImage updates the image on a recording.
// Create image if image doesn't exist on old recording but does in updated recording
// Update image if new image object is there
// Or delete image if the image has been removed
func (obj *MusicRecording) UpdateImage(db *sqlx.DB, tx *sql.Tx) (imageID IDType, err error) {
	recording, err := obj.Get(db, obj.ID)
	if err != nil {
		return imageID, fmt.Errorf("Error getting MusicRecording (%v)", err)
	}

	// If image doesn't exist on old recording but does in updated recording, create image
	if recording.Image == nil {
		if obj.Image != nil {
			err := tx.QueryRow(`
				INSERT INTO imageobject (
					cid, type, context, created_at, updated_at,
					content_url, encoding_format
				) VALUES ($1, $2, $3, $4, $5, $6, $7)
				RETURNING id`,
				obj.Image.CID, obj.Image.Type, obj.Context, time.Now(),
				time.Now(), obj.Image.ContentURL, obj.Image.EncodingFormat,
			).Scan(&imageID)

			if err != nil {
				return imageID, fmt.Errorf("Error inserting into imageobject (%v)", err)
			}

			obj.Image.ID = imageID
		}
	} else if obj.Image != nil {
		// Update image
		_, err = tx.Exec(`
			UPDATE imageobject
			SET cid=$1, type=$2, context=$3,
			updated_at=$4, content_url=$5, encoding_format=$6
			WHERE id=$7`,
			obj.Image.CID, obj.Image.Type, obj.Image.Context, time.Now(),
			obj.Image.ContentURL, obj.Image.EncodingFormat, obj.Image.ID,
		)

		if err != nil {
			return imageID, fmt.Errorf("Error updating imageobject (%v)", err)
		}

		imageID = obj.Image.ID
	} else {
		// Delete the image because it's been removed
		_, err = tx.Exec(`
			DELETE FROM person
			WHERE person.id = $1`,
			recording.Image.ID,
		)

		if err != nil {
			return imageID, fmt.Errorf("Error deleting image (%v)", err)
		}
	}

	return imageID, nil
}

// Delete a MusicRecording
// TODO - Potentially need to add CASCADES on other models
func (obj *MusicRecording) Delete(db *sqlx.DB, musicRecordingID IDType) (err error) {
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
            WHERE musicrecording.id=$1
        )`, musicRecordingID,
	)

	if err != nil {
		return fmt.Errorf("Error deleting audioobject (%v)", err)
	}

	_, err = db.Exec(`
		DELETE from musiccomposition_composer
			WHERE musiccomposition_composer.id in (
			SELECT musiccomposition_composer.id from musiccomposition_composer
			INNER JOIN musiccomposition ON musiccomposition.id = musiccomposition_composer.musiccomposition_id
			INNER JOIN musicrecording ON musicrecording.recording_of_id = musiccomposition.id
			WHERE musicrecording.id=$1
		)`, musicRecordingID,
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
        WHERE musicrecording.id=$1
    )`, musicRecordingID,
	)
	if err != nil {
		return fmt.Errorf("Error deleting musiccomposition (%v)", err)
	}

	_, err = db.Exec(
		`DELETE FROM musicrecording WHERE id=$1`, musicRecordingID,
	)
	if err != nil {
		return fmt.Errorf("Error deleting musicrecording (%v)", err)
	}
	return nil
}

// GetAudio the AudioObject for a MusicRecording
func (obj *MusicRecording) GetAudio(db *sqlx.DB, audioObjectID IDType) (*AudioObject, error) {
	query := `
    SELECT audioobject.id, audioobject.cid, audioobject.type,
        audioobject.context, audioobject.created_at, audioobject.updated_at,
        audioobject.content_url, audioobject.encoding_format
    FROM audioobject
    INNER JOIN musicrecording ON musicrecording.audio_id = audioobject.id
    WHERE musicrecording.audio_id=$1
    `

	audioobject := &AudioObject{}
	err := db.Get(audioobject, query, audioObjectID)
	if err != nil {
		return nil, fmt.Errorf("Error getting MusicRecording audio (%v)", err)
	}
	return audioobject, nil
}

// GetComposition gets the MusicComposition for a MusicRecording
func (obj *MusicRecording) GetComposition(db *sqlx.DB, musicCompositionID IDType) (*MusicComposition, error) {
	query := `
    SELECT musiccomposition.id, musiccomposition.cid, musiccomposition.type, musiccomposition.context,
    musiccomposition.created_at, musiccomposition.updated_at, musiccomposition.name, musiccomposition.iswc
    FROM musiccomposition
    INNER JOIN musicrecording ON musicrecording.recording_of_id = musiccomposition.id
    WHERE musicrecording.recording_of_id=$1
    `

	musiccomposition := &MusicComposition{}
	err := db.Get(musiccomposition, query, musicCompositionID)
	if err != nil {
		return nil, fmt.Errorf("Error getting MusicRecording composition (%v)", err)
	}
	return musiccomposition, nil
}

// GetRights gets all of the Rights objects for a MusicRecording
func (obj *MusicRecording) GetRights(db *sqlx.DB, id IDType) ([]*Right, error) {
	query := &SelectQuery{}
	query.SetQuery(`
    SELECT "right".id, "right".cid, "right".type, "right".context,
        "right".created_at, "right".updated_at, "right".percentage_shares, "right".valid_from,
        "right".valid_through, "right".person_id
    FROM "right"
    WHERE "right".musicrecording_id=$1
    `, id)

	rights := []*Right{}
	err := db.Select(&rights, query.Query(), query.Args()...)
	if err != nil {
		return nil, fmt.Errorf("Error in GetRights (%v)", err)
	}

	// Attach the party to each right object
	for _, right := range rights {
		person := &Person{}
		person, err = right.GetParty(db, *right.PartyID)
		if err != nil {
			return nil, fmt.Errorf("Error getting party (%v)", err)
		}

		right.Party = person

	}

	return rights, nil
}
