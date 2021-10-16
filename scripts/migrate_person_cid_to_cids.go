package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/consensys/ugo/pkg/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq" // postgres driver
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
		FROM person
	`
	persons := []*models.Person{}
	err = db.Select(&persons, query)
	if err != nil {
		panic(fmt.Errorf("Error getting all Persons (%v)", err))
	}
	for _, obj := range persons {
		if err != nil {
			panic(err)
		}
		query := `
		SELECT EXISTS(
			SELECT 1
			FROM person
			WHERE $1 <@ cids
		);
		`
		var exists bool
		CIDs := []string{obj.CID} // pq.Array requires an array type
		err := db.QueryRow(query, pq.Array(CIDs)).Scan(&exists)
		if err != nil {
			panic(err)
		}
		if exists {
			fmt.Println(fmt.Errorf("Skipping... CID already exists: %s", obj.CID))
			continue
		}
		fmt.Println(fmt.Errorf("Copying: %s", obj.CID))
		_, err = tx.Exec(
			`UPDATE person
			SET cids=array_append(cids, $1), updated_at=$2
			WHERE id=$3`,
			obj.CID, time.Now(), obj.ID,
		)
		if err != nil {
			panic(fmt.Errorf("Error updating person (%v)", err))
		}
	}
}
