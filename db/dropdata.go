package main

//
// import (
// 	"database/sql"
// 	"log"
//
// 	_ "github.com/lib/pq" // postgres driver
// 	migrate "github.com/rubenv/sql-migrate"
// )
//
// var db *sql.DB
//
// func teardown() {
// 	var err error
// 	db, err = sql.Open("postgres", "postgres://alexander:@localhost/ujo?sslmode=disable")
// 	if err != nil {
// 		log.Panic(err)
// 	}
//
// 	err = db.Ping()
// 	if err != nil {
// 		panic(err)
// 	}
// 	migrations := &migrate.FileMigrationSource{
// 		Dir: "db/migrations",
// 	}
// 	_, err = migrate.Exec(db, "postgres", migrations, migrate.Down)
// 	if err != nil {
// 		panic(err)
// 	}
// 	log.Printf("Database Cleared")
// }
//
// func main() {
// 	teardown()
// }
