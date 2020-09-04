package main

import (
	"log"
	"nazwa/misc"
	"nazwa/router"

	"github.com/gin-gonic/contrib/sessions"
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
	server.Use(sessions.Sessions("nazwasession", sessions.NewCookieStore([]byte("secret"))))
	server.Use(misc.NewDefaultConfig())

	server.Static("/assets", "./statics")
	server.StaticFile("/favicon.ico", "./statics/favicon.ico")

	server.LoadHTMLGlob("./templates/*")

	// Router
	// Halaman muka
	server.GET("/", router.PageHome)
	server.GET("/login", router.PageLogin)
	server.GET("/logout", router.PageLogout)
	server.GET("/create-account", router.PageCreateAccount)
	server.GET("/forgot-password", router.PageForgot)

	// Konfigurasi default admin
	server.Use(misc.NewAdminDefaultConfig())
	// Halaman Admin
	admin := server.Group("/admin")
	admin.GET("/", router.PageAdmin)
	admin.GET("/customers", router.PageAdminCustomers)
	admin.GET("/blank", router.PageAdminBlank)

	api := server.Group("/api")
	api.POST("/login", router.APIUserLogin)

	server.NoRoute(router.PageNoRoute)

	// Jalankan server
	server.Run()
}
