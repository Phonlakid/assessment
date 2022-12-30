package main

import (
	"github.com/Phonlakid/assessment/db"
)

func main() {
	db.Connect()
	db.CreateTable()
	defer db.Conn.Close()
}
