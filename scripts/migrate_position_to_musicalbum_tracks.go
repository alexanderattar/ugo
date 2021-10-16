package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/consensys/ugo/pkg/models"
	"github.com/jmoiron/sqlx"
	// postgres driver
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL must be set as environment variable")
	}
	db, err := sqlx.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}

	tx, err := db.Begin()
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

	query := `
		SELECT *
		FROM musicrecording
	`

	musicrecordings := []*models.MusicRecording{}
	err = db.Select(&musicrecordings, query)

	if err != nil {
		panic(fmt.Errorf("Error getting all MusicRecordings (%v)", err))
	}

	for _, obj := range musicrecordings {
		if err != nil {
			panic(err)
		}

		// Convert string positions to integers
		trackPosition, err := strconv.Atoi(*obj.Position)
		if err == nil {
			fmt.Println(trackPosition)
		}

		fmt.Println(fmt.Errorf(
			"Copying musicrecording position where musicrecording_id=%v, position=%v", obj.ID, trackPosition),
		)
		_, err = tx.Exec(
			`UPDATE musicalbum_tracks
			SET position=$1
			WHERE musicrecording_id=$2`,
			trackPosition, obj.ID,
		)

		if err != nil {
			panic(fmt.Errorf("Error updating musicalbum_tracks (%v)", err))
		}

	}
}
