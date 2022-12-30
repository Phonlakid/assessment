package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var Conn *sql.DB

func Connect() {

	db, err := sql.Open("postgres", "postgres://qgbloyzo:4RO9-nIJ0XNX4CArZTUb6_WJAHq6CM5O@john.db.elephantsql.com/qgbloyzo")
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	Conn = db

}

func CreateTable() {
	createTb := `
	CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);
	`
	_, err := Conn.Exec(createTb)

	if err != nil {
		log.Fatal("can't create table", err)
	}
}
