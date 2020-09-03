package main

import (
	"log"
	"nazwa/misc"
	"nazwa/router"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("Tidak ada file .env")
	}
}

func main() {
	server := gin.Default()
	server.Use(misc.NewDefaultConfig())

	server.Static("/assets", "./statics")
	server.StaticFile("/favicon.ico", "./statics/favicon.ico")

	server.LoadHTMLGlob("./templates/*")

	// Router
	// Halaman muka
	server.GET("/", router.PageHome)
	server.GET("/login", router.PageLogin)
	server.GET("/create-account", router.PageCreateAccount)
	server.GET("/forgot-password", router.PageForgot)

	// Konfigurasi default admin
	server.Use(misc.NewAdminDefaultConfig())
	// Halaman Admin
	admin := server.Group("/admin")
	admin.GET("/", router.PageAdmin)

	server.NoRoute(router.PageNoRoute)

	// Jalankan server
	server.Run()
}
