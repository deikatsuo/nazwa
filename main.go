package main

import (
	"fmt"
	"log"
	"nazwa/middleware"
	"nazwa/misc"
	"nazwa/router"
	"nazwa/router/api"
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
			setup.RunSetup(db)
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
	server.SetFuncMap(middleware.RegTmplFunc())

	// Buat session
	server.Use(sessions.Sessions("NAZWA_SESSION", sessions.NewCookieStore([]byte("secret"))))
	server.Use(middleware.NewDefaultConfig())

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
	server.GET("/create-account", router.PageCreateAccount)
	server.GET("/forgot-password", router.PageForgot)
	server.GET("/logout", router.PageLogout)
	// Halaman tidak ditemukan
	server.NoRoute(router.Page404)

	// Halaman Dashboard
	// /dashboard
	dashboard := server.Group("/dashboard")
	dashboard.Use(middleware.RoutePermission(db, e))
	// Middleware untuk mengambil pengaturan default untuk dashboard
	dashboard.Use(middleware.NewDashboardDefaultConfig(db))
	dashboard.GET("/", router.PageDashboard)
	dashboard.GET("/account", router.PageDashboardAccount(db))
	dashboard.GET("/users", router.PageDashboardUsers(db))
	dashboard.GET("/blank", router.PageDashboardBlank)

	// API
	// /api
	apis := server.Group("/api")
	apis.Use(middleware.RoutePermission(db, e))

	// V1
	// /api/v1
	v1 := apis.Group("/v1")

	// API yang diakses dari Lokal
	// /api/v1/local
	v1local := v1.Group("/local")
	v1local.POST("/login", api.UserLogin(db))
	v1local.POST("/create-account", api.UserCreate(db))

	// /api/v1/local/users
	customer := v1local.Group("/users")
	customer.GET("/list", api.UsersList(db))

	// /api/v1/local/address
	address := v1local.Group("/address")
	address.GET("/provinces", api.PlaceProvinces(db))
	address.GET("/cities/:id", api.PlaceCities(db))
	address.GET("/districts/:id", api.PlaceDistricts(db))
	address.GET("/villages/:id", api.PlaceVillages(db))

	// User API
	// /api/v1/local/user
	v1user := v1local.Group("/user")
	v1user.PATCH("/:id/update/contact", api.UserUpdateContact(db))
	v1user.DELETE("/:id/delete/email", api.UserDeleteEmail(db))
	v1user.POST("/:id/add/email", api.UserAddEmail(db))
	v1user.DELETE("/:id/delete/phone", api.UserDeletePhone(db))
	v1user.POST("/:id/add/phone", api.UserAddPhone(db))
	v1user.POST("/:id/add/address", api.UserAddAddress(db))
	v1user.DELETE("/:id/delete/address", api.UserDeleteAddress(db))

	// Jalankan server
	server.Run()
}
