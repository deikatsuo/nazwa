package main

import (
	"encoding/gob"
	"fmt"
	"nazwa/dbquery"
	"nazwa/middleware"
	"nazwa/misc"
	"nazwa/misc/validation"
	"nazwa/router"
	"nazwa/router/api"
	"nazwa/setup"
	"nazwa/wrapper"
	"net/http"
	"os"

	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/autotls"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	stats "github.com/semihalev/gin-stats"
	"golang.org/x/crypto/acme/autocert"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var log = misc.Log

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Tidak ada file .env")
	}

	// Membuat koneksi database
	log.Info("Mencoba membuat koneksi ke database...")
	db, err := sqlx.Connect(misc.SetupDBType(), misc.SetupDBSource())
	if err != nil {
		log.Info("Gagal membuat koneksi ke database ")
		log.Fatal(err)
	}

	// Set global database
	dbquery.DB = db

	gob.Register(wrapper.User{})
}

func main() {
	// Ambil argumen CLI
	if iag := len(os.Args); iag > 1 {
		arg := os.Args[1]
		switch arg {
		case "run":
			log.Info("Menjalankan server...")
			runServer()
		case "setup":
			log.Info("Menjalankan konfigurasi database...")
			setup.RunSetup()
		case "version":
			fmt.Println("Authored by", misc.AUTHOR)
			fmt.Println("Version ", misc.VERSION)
		default:
			log.Warn("Argument salah...")
		}
	}
}

func runServer() {
	// Ambil konfigurasi role
	e, err := casbin.NewEnforcer("auth_model.conf", "auth_policy.csv")

	if err != nil {
		log.Warn("Casbin enforcer fail")
		log.Fatal(err)
	}

	// gin.SetMode(gin.ReleaseMode)

	// Buat server
	server := gin.Default()

	// Redirect www ke non-www
	server.Use(middleware.RedirectWWW())

	// Kompress menggunakan gzip
	server.Use(gzip.Gzip(gzip.BestCompression))

	// Menambahkan validator date
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("date", validation.CustomValidationDate())
	}

	// Daftarkan fungsi ke template
	server.SetFuncMap(middleware.RegTmplFunc())

	// Buat session
	server.Use(sessions.Sessions("NAZWA_SESSION", sessions.NewCookieStore([]byte("deri and rika"))))
	server.Use(middleware.NewDefaultConfig)

	// Daftarkan aset statik
	// misal css, js, dan beragam file gambar
	server.Static("/assets", "./statics")
	server.Static("/file", "../data/upload")
	server.StaticFile("/favicon.ico", "./statics/favicon.ico")

	// Load file template
	server.LoadHTMLGlob("./templates/*")

	// Router
	// Halaman muka
	server.GET("/", router.PageHome)
	server.GET("/product", router.PageProduct)
	server.GET("/product/:id", router.PageProductDetail)
	server.GET("/login", router.PageLogin)
	server.GET("/create-account", router.PageCreateAccount)
	server.GET("/forgot-password", router.PageForgot)
	server.GET("/logout", router.PageLogout)
	// Halaman tidak ditemukan
	server.NoRoute(router.Page404)

	// Halaman Dashboard
	// /dashboard
	dashboard := server.Group("/dashboard")
	// Gunakan permision
	dashboard.Use(middleware.RoutePermission(e))
	// Middleware untuk mengambil pengaturan default untuk dashboard
	dashboard.Use(middleware.NewDashboardDefaultConfig())

	// PProf
	pprof.RouteRegister(dashboard, "pprof")

	// Stats
	dashboard.Use(stats.RequestStats())
	dashboard.GET("/stats", func(c *gin.Context) {
		c.JSON(http.StatusOK, stats.Report())
	})

	// dashboard.GET("/metrics", gin.WrapH(promhttp.Handler()))

	dashboard.GET("/", router.PageDashboard)
	dashboard.GET("/account", router.PageDashboardAccount)
	dashboard.GET("/users", router.PageDashboardUsers)
	dashboard.GET("/users/add", router.PageDashboardUsersAdd)
	dashboard.GET("/products", router.PageDashboardProducts)
	dashboard.GET("/products/add", router.PageDashboardProductsAdd)
	dashboard.GET("/orders", router.PageDashboardOrders)
	dashboard.GET("/orders/add", router.PageDashboardOrdersAdd)
	dashboard.GET("/orders/zones", router.PageDashboardOrdersZones)
	dashboard.GET("/blank", router.PageDashboardBlank)

	// API
	// /api
	apis := server.Group("/api")
	apis.Use(middleware.RoutePermission(e))

	// V1
	// /api/v1
	v1 := apis.Group("/v1")

	// API untuk publik
	v1public := v1.Group("/public")
	v1public.GET("/product/all", api.ProductShowAll)

	// API yang diakses dari Lokal
	// /api/v1/local
	v1local := v1.Group("/local")
	v1local.POST("/login", api.UserLogin)
	//v1local.POST("/create-account", api.UserCreate(db))

	// /api/v1/local/address
	v1address := v1local.Group("/address")
	v1address.GET("/provinces", api.PlaceProvinces)
	v1address.GET("/cities/:id", api.PlaceCities)
	v1address.GET("/districts/:id", api.PlaceDistricts)
	v1address.GET("/villages/:id", api.PlaceVillages)

	// /api/v1/local/zone
	v1zone := v1local.Group("/zone")
	v1zone.GET("/list", api.ZoneGetList)

	// /api/v1/local/product
	v1product := v1local.Group("/product")
	v1product.GET("/id/:id", api.ProductShowByID)
	v1product.GET("/list/:limit", api.ProductShowList)
	v1product.GET("/all", api.ProductShowAll)
	v1product.POST("/add", api.ProductCreate)

	// /api/v1/local/product/edit
	v1pEdit := v1product.Group("/edit")
	v1pEdit.POST("/:id/add/credit_price", api.ProductAddCreditPrice)
	v1pEdit.DELETE("/:id/delete/credit_price", api.ProductDeleteCreditPrice)

	// /api/v1/local/product/search
	v1pSearch := v1product.Group("/search")
	v1pSearch.GET("/name/:limit", api.ProductSearchByName)

	// /api/v1/local/order
	v1order := v1local.Group("/order")
	v1order.GET("/id/:id", api.OrderShowByID)
	v1order.GET("/list/:limit", api.OrderShowList)
	v1order.POST("/create", api.OrderCreate)
	v1order.GET("/substitute/ric", api.OrderSubstituteByRicShow)

	// /api/v1/local/order/edit
	v1oEdit := v1order.Group("/edit")
	v1oEdit.DELETE("/:id/delete", api.OrderDeleteByID)

	// User API
	// /api/v1/local/user
	v1user := v1local.Group("/user")
	v1user.POST("/create", api.UserCreate)
	v1user.GET("/list/:limit", api.UserShowList)
	v1user.GET("/id/:id", api.UserShowByID)
	v1user.GET("/address/:id", api.UserShowAddressByUserID)

	// User API edit
	// /api/v1/local/user/edit
	v1uEdit := v1user.Group("/edit")
	v1uEdit.PATCH("/:id/update/contact", api.UserUpdateContact)
	v1uEdit.PATCH("/:id/update/role", api.UserUpdateRole)
	v1uEdit.DELETE("/:id/delete/email", api.UserDeleteEmail)
	v1uEdit.POST("/:id/add/email", api.UserAddEmail)
	v1uEdit.DELETE("/:id/delete/phone", api.UserDeletePhone)
	v1uEdit.POST("/:id/add/phone", api.UserAddPhone)
	v1uEdit.POST("/:id/add/address", api.UserAddAddress)
	v1uEdit.DELETE("/:id/delete/address", api.UserDeleteAddress)

	// User API search/pencarian
	// /api/v1/local/user/search
	v1uSearch := v1user.Group("/search")
	v1uSearch.GET("/ric/:limit", api.UserSearchByNIK)
	v1uSearch.GET("/collector/:limit", api.UserSearchByNameType("2"))
	v1uSearch.GET("/surveyor/:limit", api.UserSearchByNameType("4"))
	v1uSearch.GET("/sales/:limit", api.UserSearchByNameType("5"))

	port := ":" + misc.GetEnv("PORT", "8080").(string)

	if misc.GetEnv("REMOTE", "false") == "true" {
		var hostname string
		if misc.GetEnv("HOSTNAME", "") != "" {
			hostname = misc.GetEnvND("HOSTNAME")
		}
		cert := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(hostname, "www."+hostname),
			Cache:      autocert.DirCache("../data/cert_cache"),
		}

		// Jalankan server dalam mode aman
		log.Fatal(autotls.RunWithManager(server, &cert))
	} else {
		// Jalankan server dalam mode tidak aman
		server.Run(port)
	}
}
