package main

import (
	"fmt"
	"log"
	"nazwa/middleware"
	"nazwa/misc"
	"nazwa/router"
	"nazwa/setup"
	"os"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("Tidak ada file .env")
	}
}

func main() {
	// Membuat koneksi database
	fmt.Println("Mencoba membuat koneksi ke database...")
	db, err := sqlx.Connect(misc.SetupDBType(), misc.SetupDBSource())
	if err != nil {
		fmt.Println("Gagal membuat koneksi ke database ")
		fmt.Println(err)
		os.Exit(1)
	}
	// Ambil argumen CLI
	if iag := len(os.Args); iag > 1 {
		arg := os.Args[1]
		switch arg {
		case "run":
			fmt.Println("Menjalankan server...")
			runServer(db)
		case "setup":
			fmt.Println("Menjalankan konfigurasi database...")
			setup.RunSetup(db, false)
		case "setup-reset":
			fmt.Println("PERINGATAN: Ini akan menghapus semua data di database!")
			var lanjut string
			fmt.Print("Lanjutkan? [y/N]")
			fmt.Scanf("%s", &lanjut)
			if lanjut == "Y" || lanjut == "y" {
				fmt.Println("Menjalankan konfigurasi dan mereset database...")
				setup.RunSetup(db, true)
			}
		case "version":
			fmt.Println("Authored by", misc.AUTHOR)
			fmt.Println("Version ", misc.VERSION)
		default:
			fmt.Println("Argument salah...")
		}
	}
}

func runServer(db *sqlx.DB) {
	// Ambil konfigurasi role
	e, err := casbin.NewEnforcer("auth_model.conf", "auth_policy.csv")
	if err != nil {
		log.Fatal(err)
	}

	// Buat server
	server := gin.Default()

	// Daftarkan fungsi ke template
	server.SetFuncMap(misc.RegTmplFunc())

	// Buat session
	server.Use(sessions.Sessions("NAZWA_SESSION", sessions.NewCookieStore([]byte("secret"))))
	server.Use(misc.NewDefaultConfig())

	// Periksa user role
	server.Use(middleware.RoutePermission(db, e))

	// Daftarkan aset statik
	// misal css, js, dan beragam file gambar
	server.Static("/assets", "./statics")
	server.Static("/file", "./upload")
	server.StaticFile("/favicon.ico", "./statics/favicon.ico")

	// Load file template
	server.LoadHTMLGlob("./templates/*")

	// Router
	// Halaman muka
	server.GET("/", router.PageHome)
	server.GET("/login", router.PageLogin)
	server.GET("/logout", router.PageLogout)
	server.GET("/create-account", router.PageCreateAccount)
	server.GET("/forgot-password", router.PageForgot)
	// Halaman tidak ditemukan
	server.NoRoute(router.Page404)

	// Halaman Dashboard
	dashboard := server.Group("/dashboard")
	// Middleware untuk mengambil pengaturan default untuk dashboard
	dashboard.Use(misc.NewDashboardDefaultConfig(db))
	dashboard.GET("/", router.PageDashboard)
	dashboard.GET("/customers", router.PageDashboardCustomers)
	dashboard.GET("/blank", router.PageDashboardBlank)

	// API
	api := server.Group("/api")

	// V1
	v1 := api.Group("/v1")

	v1local := v1.Group("/local")
	v1local.POST("/login", router.APIUserLogin(db))
	v1local.POST("/create-account", router.APIUserCreate)

	// Jalankan server
	server.Run()
}
