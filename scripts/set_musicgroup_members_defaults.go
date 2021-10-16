package main

import (
	"fmt"
	"log"
	"os"

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
		FROM musicgroup 
	`

	musicgroups := []*models.MusicGroup{}
	err = db.Select(&musicgroups, query)

	if err != nil {
		panic(fmt.Errorf("Error getting all MusicGroups (%v)", err))
	}

	for _, obj := range musicgroups {
		if err != nil {
			panic(err)
		}

		fmt.Println(fmt.Errorf(
			"Setting defaults for musicgroup_members where musicgroup_id=%v", obj.ID),
		)

		_, err = tx.Exec(
			`UPDATE musicgroup_members
			SET percentage_shares=$1, musicgroup_admin=$2
			WHERE musicgroup_id=$3`,
			100, true, obj.ID,
		)

		if err != nil {
			panic(fmt.Errorf("Error updating musicalbum_tracks (%v)", err))
		}
	}
}
