package setup

import (
	"fmt"
	"nazwa/dbquery"
	"nazwa/misc"
	"os"

	"github.com/jmoiron/sqlx"
)

// RunSetup ...
// Menjalankan setup
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

	fmt.Println("Membuat tabel...")
	if misc.Migration("down") {
		fmt.Println("Mereset tabel ('-')")
	}
	if misc.Migration("up") {
		fmt.Println("Membuat tabel ('-')")
	}

	fmt.Println("Setup user admin...")
	setupUserAdmin(db)
}

// Membuat user admin baru
func setupUserAdmin(db *sqlx.DB) error {
	user := dbquery.NewUser()

	var firstname, lastname, username, password, gender string

	fmt.Print("\033[H\033[2J")
	fmt.Print("Nama depan: ")
	fmt.Scanf("%s", &firstname)
	fmt.Print("Nama belakang: ")
	fmt.Scanf("%s", &lastname)
	fmt.Println("Jenis kelamin [m/f]")
	fmt.Scanf("%s", &gender)
	fmt.Print("Username: ")
	fmt.Scanf("%s", &username)
	fmt.Print("Password: ")
	fmt.Scanf("%s", &password)

	var uid int
	err := user.SetFirstName(firstname).
		SetLastName(lastname).
		SetUserName(username).
		SetPassword(password).
		SetGender(gender).
		ReturnID(&uid).
		Save(db)

	if err != nil {
		return err
	}
	return nil
}
