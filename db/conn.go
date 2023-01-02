package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var Conn *sql.DB

func Connect() {

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	Conn = db
	CreateTable(Conn)
}

func CreateTable(Conn *sql.DB) {
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
