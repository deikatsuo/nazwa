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
func Migration(s string) bool {
	success := true
	url := SetupMigrationURL()
	m, err := migrate.New("file://setup/migration", url)

	if err != nil {
		log.Fatal(err)
	}
	if s == "up" {
		if err := m.Up(); err != nil {
			success = false
			fmt.Println("Migration status: ", err)
		}
	} else if s == "down" {
		if err := m.Down(); err != nil {
			success = false
			fmt.Println("Migration status: ", err)
		}
	}
	return success
}
