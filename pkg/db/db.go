package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

// NewDB sets up a database connection
func NewDB(connectURI string) (*sqlx.DB, error) {

	os.Getenv("DB_DRIVER")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	connectURI = fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		dbHost, dbPort, dbUser, dbName, dbPassword,
	)

	fmt.Println("++++++++++++")
	fmt.Println(connectURI)
	fmt.Println("++++++++++++")

	db, err := sqlx.Open("postgres", connectURI)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
