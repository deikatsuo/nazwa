package setup

import (
	"fmt"
	"nazwa/dbquery"
	"nazwa/misc"
	"os"

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
		fmt.Println("Gagal membuat koneksi ke database ")
		fmt.Println(err)
		os.Exit(1)
	}

	misc.Migration("up")

	fmt.Println("Setup user admin...")
	setupUserAdmin(db)
}

// Membuat user admin baru
func setupUserAdmin(db *sqlx.DB) {
	user := dbquery.NewUser()

	user.FirstName("deri").
		LastName("herdianto").
		UserName("deikatsuo").
		Password("love100").
		Gender("m")

	//tx := db.MustBegin()
}
