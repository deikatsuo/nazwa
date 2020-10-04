package setup

import (
	"fmt"
	"log"
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/misc/validation"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gocarina/gocsv"
	"github.com/jmoiron/sqlx"
)

// RunSetup menjalankan setup servers
// Menjalankan setup
func RunSetup(db *sqlx.DB) {
	var reset bool

	// Tanya user apakah ingin mereset database
	var doreset string
	fmt.Println()
	fmt.Println("Apakah anda ingin melakukan reset DATABASE? [y/N]")
	fmt.Scanf("%s", &doreset)
	if doreset == "y" || doreset == "Y" {
		reset = true
	}

	// jika ya
	// maka hapus tabel (reset) pada database
	if reset {
		fmt.Println("PERINGATAN: Ini akan menghapus semua data di database!")
		var lanjut string
		fmt.Println("Lanjutkan? [y/N]")
		fmt.Scanf("%s", &lanjut)
		if lanjut == "y" || lanjut == "Y" {
			fmt.Println("Sedang mencoba menghapus semua tabel...")
			if misc.Migration("down") {
				fmt.Println("Tabel berhasil dihapus ('-')")
			}
		}
	}

	// Upgrade tabel, atau buat baru jika belum ada
	if reset {
		fmt.Println("Mencoba kembali memulihkan tabel yang telah dihapus")
	} else {
		fmt.Println("Sedang memutakhirkan tabel")
	}
	if misc.Migration("up") {
		fmt.Println("Tabel berhasil dibuat ('-')")
	}

	fmt.Println()
	fmt.Println("Setup tabel di database selesai")
	fmt.Println()

	// Tanya user apakah ingin melakukan
	// konfigurasi daerah
	var setdaerah string
	fmt.Println("Apakah anda ingin melakukan setup daerah? [y/N]")
	fmt.Scanf("%s", &setdaerah)
	if setdaerah == "y" || setdaerah == "Y" {
		if err := setupDaerah(db); err != nil {
			fmt.Println("Terjadi kesalahan saat konfigurasi daerah")
			log.Fatal(err)
		}
		fmt.Println()
		fmt.Println("Konfigurasi daerah selesai!")
	}

	var buatAdmin string
	fmt.Print("Buat user admin? [y/N]")
	fmt.Scanf("%s", &buatAdmin)
	if buatAdmin == "y" || buatAdmin == "Y" {
		// Lakukan pendaftaran admin baru
		fmt.Println("Setup user admin...")
		if err := setupUserAdmin(db); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	fmt.Println()
	fmt.Println("Setup selesai")

}

// Membuat user admin baru
func setupUserAdmin(db *sqlx.DB) error {
	user := dbquery.NewUser()

	// Variabel untuk menyimpan hasil input
	var input createAdminInput

	fmt.Print("\033[H\033[2J")
	fmt.Print("Nama depan: ")
	fmt.Scanf("%s", &input.Firstname)
	fmt.Print("Nama belakang: ")
	fmt.Scanf("%s", &input.Lastname)
	fmt.Println("Jenis kelamin [m/f]")
	fmt.Scanf("%s", &input.Gender)
	fmt.Print("Username: ")
	fmt.Scanf("%s", &input.Username)
	fmt.Print("Nomor HP: ")
	fmt.Scanf("%s", &input.Phone)
	fmt.Print("Password: ")
	fmt.Scanf("%s", &input.Password)
	fmt.Print("Email: ")
	fmt.Scanf("%s", &input.Email)

	// Validasi hasil input dari user
	validate := validator.New()
	if err := validate.Struct(&input); err != nil {
		erbar := validation.SimpleValErr(err)
		fmt.Println(erbar)

		// Looping jika input tidak benar
		var again string
		fmt.Print("Ulangi lagi? [y/N]?")
		fmt.Scanf("%s", &again)
		if again == "y" || again == "Y" {
			setupUserAdmin(db)
			return nil
		}
		os.Exit(1)
	}

	// Variabel untuk menyimpan ID dari insert terakhir
	var uid int

	err := user.SetFirstName(input.Firstname).
		SetLastName(input.Lastname).
		SetUserName(input.Username).
		SetPhone(input.Phone).
		SetEmail(input.Email).
		SetPassword(input.Password).
		SetGender(input.Gender).
		SetRole(dbquery.RoleAdmin).
		ReturnID(&uid).
		Save(db)
	if err != nil {
		return err
	}

	return nil
}

// createAdminInput user untuk registrasi admin setelah
// melakukan setup
type createAdminInput struct {
	Firstname string `validate:"required,alpha,min=3,max=25"`
	Lastname  string `validate:"alpha,min=1,max=25"`
	Username  string `validate:"alphanum,min=4,max=25"`
	Password  string `validate:"alphanumunicode,min=8,max=25"`
	Gender    string `validate:"required,oneof=m f"`
	Phone     string `validate:"numeric,min=6,max=15"`
	Email     string `validate:"email"`
}

// Daerah - struk untuk menyimpan data daerah
type Daerah struct {
	ID        string  `csv:"Code" db:"id"`
	Parent    int     `csv:"Parent" db:"parent"`
	Name      string  `csv:"Name" db:"name"`
	Latitude  float32 `csv:"-"`
	Longitude float32 `csv:"-"`
	Postal    string  `csv:"-"`
}

func setupDaerah(db *sqlx.DB) error {
	tx := db.MustBegin()
	fmt.Println("Membuat data provinsi")
	if provinces, err := openData("provinces.csv"); err != nil {
		return err
	} else {
		query := `INSERT INTO "province" (id, parent, name) VALUES (:id, :parent, :name)`
		if _, err := tx.NamedExec(query, provinces); err != nil {
			log.Print("ERRSETUP-20")
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		log.Print("ERRSETUP-21")
		return err
	}
	return nil
}

func openData(d string) ([]Daerah, error) {
	file, err := os.Open(fmt.Sprintf("./setup/%s", d))
	//file, err := os.OpenFile("./setup/villages.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return []Daerah{}, err
	}
	defer file.Close()

	data := []Daerah{}

	if err := gocsv.UnmarshalFile(file, &data); err != nil {
		log.Print("ERRSETUP-22")
		return []Daerah{}, err
	}

	return data, nil
}
