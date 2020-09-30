package main

import (
	"fmt"
	"log"
	"nazwa/misc"
	"nazwa/router"
	"nazwa/setup"
	"os"

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
	if iag := len(os.Args); iag > 1 {
		arg := os.Args[1]
		switch arg {
		case "run":
			fmt.Println("Menjalankan server...")
			runServer()
		case "setup":
			fmt.Println("Menjalankan setup...")
			setup.RunSetup()
		case "version":
			fmt.Println("Authored by", misc.AUTHOR)
			fmt.Println("Version ", misc.VERSION)
		default:
			fmt.Println("Argument salah...")
		}
	}
}

func runServer() {
	// Membuat koneksi database
	fmt.Println("Mencoba membuat koneksi ke database...")
	db, err := sqlx.Connect(misc.SetupDBType(), misc.SetupDBSource())
	if err != nil {
		fmt.Println("Gagal membuat koneksi ke database ")
		fmt.Println(err)
		os.Exit(1)
	}

	server := gin.Default()
	server.Use(sessions.Sessions("nazwasession", sessions.NewCookieStore([]byte("secret"))))
	server.Use(misc.NewDefaultConfig())

	server.Static("/assets", "./statics")
	server.Static("/file", "./upload")
	server.StaticFile("/favicon.ico", "./statics/favicon.ico")

	server.LoadHTMLGlob("./templates/*")

	// Router
	// Halaman muka
	server.GET("/", router.PageHome)
	server.GET("/login", router.PageLogin)
	server.GET("/logout", router.PageLogout)
	server.GET("/create-account", router.PageCreateAccount)
	server.GET("/forgot-password", router.PageForgot)

	// Konfigurasi default dashboard
	server.Use(misc.NewDashboardDefaultConfig())
	// Halaman Dashboard
	dashboard := server.Group("/dashboard")
	dashboard.GET("/", router.PageDashboard)
	dashboard.GET("/customers", router.PageDashboardCustomers)
	dashboard.GET("/blank", router.PageDashboardBlank)

	api := server.Group("/api")
	api.POST("/login", router.APIUserLogin(db))
	api.POST("/create-account", router.APIUserCreate)

	server.NoRoute(router.PageNoRoute)

	// Jalankan server
	server.Run()
}
