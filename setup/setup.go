package setup

import (
	"fmt"
	"nazwa/dbquery"
	"nazwa/misc"
	"os"

	"github.com/go-playground/validator/v10"
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
	if err := setupUserAdmin(db); err != nil {
		fmt.Println(err)
	}
}

// Membuat user admin baru
func setupUserAdmin(db *sqlx.DB) error {
	user := dbquery.NewUser()

	var input adminInput

	fmt.Print("\033[H\033[2J")
	fmt.Print("Nama depan: ")
	fmt.Scanf("%s", &input.Firstname)
	fmt.Print("Nama belakang: ")
	fmt.Scanf("%s", &input.Lastname)
	fmt.Println("Jenis kelamin [m/f]")
	fmt.Scanf("%s", &input.Gender)
	fmt.Print("Username: ")
	fmt.Scanf("%s", &input.Username)
	fmt.Print("Password: ")
	fmt.Scanf("%s", &input.Password)

	validate := validator.New()
	if err := validate.Struct(&input); err != nil {
		return err
	}

	var uid int
	err := user.SetFirstName(input.Firstname).
		SetLastName(input.Lastname).
		SetUserName(input.Username).
		SetPassword(input.Password).
		SetGender(input.Gender).
		ReturnID(&uid).
		Save(db)

	if err != nil {
		return err
	}
	return nil
}

type adminInput struct {
	Firstname string `validate:"required,alpha,min=3,max=25"`
	Lastname  string `validate:"alpha,min=1,max=25"`
	Username  string `validate:"alphanum,min=4,max=25"`
	Password  string `validate:"alphanumunicode,min=8,max=25"`
	Gender    string `validate:"required,oneof=m f"`
}
