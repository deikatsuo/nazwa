package setup

import (
	"fmt"
	"io/ioutil"
	"nazwa/misc"

	"github.com/jmoiron/sqlx"
)

func RunSetup() {
	createTables()
}

func createTables() {
	// Membuat koneksi database
	fmt.Println("Mencoba membuat koneksi ke database...")
	db, err := sqlx.Connect(misc.SetupDBType(), misc.SetupDBSource())
	if err != nil {
		fmt.Println("Gagal membuat koneksi ke database ", err)
	}

	// Load tables.sql
	fmt.Println("Membaca file tables.sql")
	schema, err := ioutil.ReadFile("./setup/tables.sql")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Menjalankan sql query")
	db.MustExec(string(schema))

	fmt.Println("Selesai...")
}
