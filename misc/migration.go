package misc

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Migration ...
// Fungsi untuk melakukan migrasi
// dengan parameter "up" dan "down"
func Migration(s string) {
	fail := false
	url := SetupMigrationURL()
	m, err := migrate.New("file://setup/migration", url)
	fmt.Println("Persiapan untuk menjalankan migration...")
	if err != nil {
		log.Fatal(err)
	}
	if s == "up" {
		if err := m.Up(); err != nil {
			fail = true
			fmt.Println("Mencoba melakukan upgrade")
			fmt.Println("Migration status: ", err)
		}
	} else if s == "down" {
		if err := m.Down(); err != nil {
			fail = true
			fmt.Println("Mencoba melakukan downgrade")
			fmt.Println("Migration status: ", err)
		}
	}
	if !fail {
		fmt.Println("Berhasil melakukan migration ('-')")
	}
}
