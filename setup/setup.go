package setup

import (
	"fmt"
	"log"
	"nazwa/dbquery"
	"nazwa/misc"
	"nazwa/misc/validation"
	"nazwa/wrapper"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
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
		fmt.Println("PERINGATAN: Ini akan menghapus semua data user pada server dan database!")
		var lanjut string
		fmt.Println("Lanjutkan? [y/N]")
		fmt.Scanf("%s", &lanjut)
		if lanjut == "y" || lanjut == "Y" {
			var pwd string
			fmt.Println("Silahkan masukan kata sandi")
			fmt.Scanf("%s", &pwd)
			if pwd == misc.GetEnvND("DEV_PWD") {
				fmt.Println("Sedang mencoba menghapus semua tabel...")
				if misc.Migration("down") {
					fmt.Println("Tabel berhasil dihapus ('-')")
				}

				fmt.Println("Menghaspus data user dan produk")
				purgeUser := os.RemoveAll("../data/upload")
				if purgeUser != nil {
					log.Println("ERROR: setup.go RunSetup() gagal menghapus data user dan produk")
					log.Println(purgeUser)
				} else {
					// Membuat folder untuk menyimpan data user
					fmt.Println("Membuat kembali data user dan produk")

					if err := os.Mkdir("../data", 0755); err != nil {
						log.Println("ERROR: setup.go RunSetup() gagal membuat direktori data")
						log.Println(err)
					}

					if err := os.Mkdir("../data/upload", 0755); err != nil {
						log.Println("ERROR: setup.go RunSetup() gagal membuat direktori /upload")
						log.Println(err)
					}

					if err := os.Mkdir("../data/upload/profile", 0755); err != nil {
						log.Println("ERROR: setup.go RunSetup() gagal membuat direktori /upload/profile")
						log.Println(err)
					}

					if err := os.Mkdir("../data/upload/profile/thumbnail", 0755); err != nil {
						log.Println("ERROR: setup.go RunSetup() gagal membuat direktori /upload/profile/thumbnail")
						log.Println(err)
					}

					if err := os.Mkdir("../data/upload/product", 0755); err != nil {
						log.Println("ERROR: setup.go RunSetup() gagal membuat direktori /upload/product")
						log.Println(err)
					}

					if err := os.Mkdir("../data/upload/product/thumbnail", 0755); err != nil {
						log.Println("ERROR: setup.go RunSetup() gagal membuat direktori /upload/product/thumbnail")
						log.Println(err)
					}

					// Copy file male.png
					if err := misc.CopyFile("../data/upload/profile/male.png", "./setup/male.png"); err != nil {
						log.Println("ERROR: setup.go RunSetup() gagal meng copy file male.png")
						log.Println(err)
					}

					// Copy file thumb.male.png
					if err := misc.CopyFile("../data/upload/profile/thumbnail/male.png", "./setup/thumb.male.png"); err != nil {
						log.Println("ERROR: setup.go RunSetup() gagal meng copy file thumb.male.png")
						log.Println(err)
					}

					// Copy file female.png
					if err := misc.CopyFile("../data/upload/profile/female.png", "./setup/female.png"); err != nil {
						log.Println("ERROR: setup.go RunSetup() gagal meng copy file user")
						log.Println(err)
					}

					// Copy file thumb.female.png
					if err := misc.CopyFile("../data/upload/profile/thumbnail/female.png", "./setup/thumb.female.png"); err != nil {
						log.Println("ERROR: setup.go RunSetup() gagal meng copy file user")
						log.Println(err)
					}

					// Copy file no-photo.png
					if err := misc.CopyFile("../data/upload/product/no-photo.png", "./setup/no-photo.png"); err != nil {
						log.Println("ERROR: setup.go RunSetup() gagal meng copy file no-photo.png")
						log.Println(err)
					}

					// Copy file thumb.no-photo.png
					if err := misc.CopyFile("../data/upload/product/thumbnail/no-photo.png", "./setup/thumb.no-photo.png"); err != nil {
						log.Println("ERROR: setup.go RunSetup() gagal meng copy thumb.no-photo.png")
						log.Println(err)
					}
				}

			} else {
				fmt.Println("Kata sandi tidak sesuai")
				os.Exit(1)
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
	user := dbquery.UserNew()

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
	fmt.Print("Nomor KTP: ")
	fmt.Scanf("%s", &input.RIC)
	fmt.Print("Nomor KK: ")
	fmt.Scanf("%s", &input.FamilyCard)
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
		SetRIC(input.RIC).
		SetPhone(input.Phone).
		SetEmail(input.Email).
		SetPassword(input.Password).
		SetGender(input.Gender).
		SetRole(wrapper.UserRoleDev).
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
	Firstname  string `validate:"required,alpha,min=3,max=25"`
	Lastname   string `validate:"alpha,min=1,max=25"`
	Username   string `validate:"alphanum,min=4,max=25"`
	Password   string `validate:"alphanumunicode,min=8,max=25"`
	Gender     string `validate:"required,oneof=m f"`
	RIC        string `validate:"numeric,min=16,max=16"`
	FamilyCard string `validate:"numeric,min=16,max=16"`
	Phone      string `validate:"numeric,min=6,max=15"`
	Email      string `validate:"email"`
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
	var err error
	var query string
	var provinces []Daerah
	var city []Daerah
	var district []Daerah
	var village []Daerah

	fmt.Println()
	// PROVINSI
	fmt.Println(">>Membuat data provinsi")
	// Load data provinsi dari file csv
	provinces, err = openData("provinces.csv")
	if err != nil {
		log.Print("ERRSETUP-19")
		return err
	}
	// Masukan data provinsi ke database
	query = `INSERT INTO "province" (id, parent, name) VALUES (:id, :parent, :name)`
	if _, err := tx.NamedExec(query, provinces); err != nil {
		log.Print("ERRSETUP-20")
		return err
	}

	// KOTA/KABUPATEN
	fmt.Println(">>Membuat data kota/kabupaten")
	// Load data kota/kabupaten dari csv
	city, err = openData("cities.csv")
	if err != nil {
		log.Println("ERRSETUP-18")
		return err
	}
	// Masukan data kota/kabupaten ke database
	query = `INSERT INTO "city" (id, parent, name) VALUES (:id, :parent, :name)`
	if _, err := tx.NamedExec(query, city); err != nil {
		log.Print("ERRSETUP-17")
		return err
	}

	// DISTRIK/KECAMATAN
	fmt.Println(">>Membuat data distrik/kecamatan")
	// Load data distrik/kecamatan dari csv
	district, err = openData("sub-districts.csv")
	if err != nil {
		log.Println("ERRSETUP-16")
		return err
	}
	// Masukan data distrik/kecamatan ke database
	query = `INSERT INTO "district" (id, parent, name) VALUES (:id, :parent, :name)`
	if _, err := tx.NamedExec(query, district); err != nil {
		log.Print("ERRSETUP-15")
		return err
	}

	// KELURAHAN/DESA
	fmt.Println(">>Membuat data kelurahan/desa")
	// Load data kelurahan/desa dari csv
	village, err = openData("villages.csv")
	if err != nil {
		log.Println("ERRSETUP-14")
		return err
	}
	// Karena data kelurahan terlalu besar
	// maka harus di split per 20000 insert
	split := 20000
	start := 0
	vilen := len(village)
	query = `INSERT INTO "village" (id, parent, name) VALUES (:id, :parent, :name)`
	count := 5
	bar := pb.Simple.Start(count)
	for {
		bar.Increment()
		time.Sleep(time.Millisecond)
		if (start + split) < vilen {
			if _, err := tx.NamedExec(query, village[start:start+split]); err != nil {
				log.Print("ERRSETUP-13")
				return err
			}
			start = start + split
		} else {
			if _, err := tx.NamedExec(query, village[start:]); err != nil {
				log.Print("ERRSETUP-13")
				return err
			}
			break
		}
	}
	bar.Finish()

	// Comit
	if err := tx.Commit(); err != nil {
		log.Print("ERRSETUP-21")
		return err
	}
	fmt.Println("Selesai membuat data daerah")
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
