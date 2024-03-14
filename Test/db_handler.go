package controllers

import (
	"database/sql"
	"log"
)

func connect() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/db_latihan_pbp?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
