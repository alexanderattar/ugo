package main

import (
	"database/sql"
	"log"
	"os"

	testfixtures "gopkg.in/testfixtures.v2"

	_ "github.com/lib/pq" // postgres driver
	migrate "github.com/rubenv/sql-migrate"
)

var (
	db       *sql.DB
	fixtures *testfixtures.Context
)

func setup() {
	username := os.Getenv("DATABASE_USERNAME")
	password := os.Getenv("DATABASE_PASSWORD")

	var err error
	db, err = sql.Open("postgres", "postgres://"+username+":"+password+"@localhost/ujo?sslmode=disable")
	if err != nil {
		log.Panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	migrations := &migrate.FileMigrationSource{
		Dir: "db/migrations",
	}

	// Make sure database is clear
	_, err = migrate.Exec(db, "postgres", migrations, migrate.Down)
	if err != nil {
		panic(err)
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		panic(err)
	}

	log.Printf("Database Initialized. Applied %d migrations!\n", n)
}

func main() {
	setup()
	log.Print("Loading fixtures")
	// TODO - Disable SkipDatabaseNameCheck in production
	// Open connection with the test database.
	// Do NOT import fixtures in a production database!
	// Existing data would be deleted
	testfixtures.SkipDatabaseNameCheck(true)
	db, err := sql.Open("postgres", "postgres://docker:docker@localhost/ujo?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	// creating the context that hold the fixtures
	// see about all compatible databases in this page below
	fixtures, err = testfixtures.NewFolder(db, &testfixtures.PostgreSQL{}, "db/fixtures")
	if err != nil {
		log.Fatal(err)
	}

	err = fixtures.Load()
	if err != nil {
		log.Fatal(err)
	}
}
