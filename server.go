package main

import (
	"github.com/Phonlakid/assessment/db"
)

func main() {
	db.Connect()

	defer db.Conn.Close()
}
