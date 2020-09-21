package main

import (
	"fmt"
	"io/ioutil"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// Membuat koneksi database
	db, err := sqlx.Connect("postgres", "user=postgres dbname=nazwa sslmode=disable")
	if err != nil {
		fmt.Println("Gagal membuat koneksi ke database ", err)
	}

	// Load tables.sql
	schema, err := ioutil.ReadFile("./setup/tables.sql")
	if err != nil {
		fmt.Println(err)
	}

	db.MustExec(string(schema))
}
