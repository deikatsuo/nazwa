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
	server.GET("/", router.Homepage)
	server.GET("/login", router.Login)
	server.GET("/forgot-password", router.Forgot)
	server.NoRoute(router.NoRoute)

	// Jalankan server
	server.Run()
}
