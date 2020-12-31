package main

import (
	"fmt"
	"log"
	"nazwa/middleware"
	"nazwa/misc"
	"nazwa/misc/validation"
	"nazwa/router"
	"nazwa/router/api"
	"nazwa/setup"
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
	server.Use(middleware.NewDefaultConfig())

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
	dashboard.Use(middleware.RoutePermission(db, e))
	// Middleware untuk mengambil pengaturan default untuk dashboard
	dashboard.Use(middleware.NewDashboardDefaultConfig(db))

	// PProf
	pprof.RouteRegister(dashboard, "pprof")

	// Stats
	dashboard.Use(stats.RequestStats())
	dashboard.GET("/stats", func(c *gin.Context) {
		c.JSON(http.StatusOK, stats.Report())
	})

	// dashboard.GET("/metrics", gin.WrapH(promhttp.Handler()))

	dashboard.GET("/", router.PageDashboard)
	dashboard.GET("/account", router.PageDashboardAccount(db))
	dashboard.GET("/users", router.PageDashboardUsers(db))
	dashboard.GET("/users/add", router.PageDashboardUsersAdd(db))
	dashboard.GET("/products", router.PageDashboardProducts(db))
	dashboard.GET("/products/add", router.PageDashboardProductsAdd(db))
	dashboard.GET("/orders", router.PageDashboardOrders(db))
	dashboard.GET("/orders/add", router.PageDashboardOrdersAdd(db))
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
	//v1local.POST("/create-account", api.UserCreate(db))

	// /api/v1/local/address
	v1address := v1local.Group("/address")
	v1address.GET("/provinces", api.PlaceProvinces(db))
	v1address.GET("/cities/:id", api.PlaceCities(db))
	v1address.GET("/districts/:id", api.PlaceDistricts(db))
	v1address.GET("/villages/:id", api.PlaceVillages(db))

	// /api/v1/local/product
	v1product := v1local.Group("/product")
	v1product.GET("/id/:id", api.ProductShowByID(db))
	v1product.GET("/list/:limit", api.ProductShowList(db))
	v1product.GET("/all", api.ProductShowAll(db))
	v1product.POST("/add", api.ProductCreate(db))

	// /api/v1/local/product/search
	v1pSearch := v1product.Group("/search")
	v1pSearch.GET("/name/:limit", api.ProductSearchByName(db))

	// /api/v1/local/order
	v1order := v1local.Group("/order")
	v1order.GET("/id/:id", api.OrderShowByID(db))
	v1order.GET("/list/:limit", api.OrderShowList(db))
	v1order.POST("/create", api.OrderCreate(db))

	// User API
	// /api/v1/local/user
	v1user := v1local.Group("/user")
	v1user.POST("/create", api.UserCreate(db))
	v1user.GET("/list/:limit", api.UserShowList(db))
	v1user.GET("/id/:id", api.UserShowByID(db))
	v1user.GET("/address/:id", api.UserShowAddressByUserID(db))

	// User API search/pencarian
	// /api/v1/local/user/search
	v1uSearch := v1user.Group("/search")
	v1uSearch.GET("/ric/:limit", api.UserSearchByNIK(db))
	v1uSearch.GET("/collector/:limit", api.UserSearchByNameType(db, "2"))
	v1uSearch.GET("/surveyor/:limit", api.UserSearchByNameType(db, "4"))
	v1uSearch.GET("/sales/:limit", api.UserSearchByNameType(db, "5"))

	// User API edit
	// /api/v1/local/user/edit
	v1uEdit := v1user.Group("/edit")
	v1uEdit.PATCH("/:id/update/contact", api.UserUpdateContact(db))
	v1uEdit.PATCH("/:id/update/role", api.UserUpdateRole(db))
	v1uEdit.DELETE("/:id/delete/email", api.UserDeleteEmail(db))
	v1uEdit.POST("/:id/add/email", api.UserAddEmail(db))
	v1uEdit.DELETE("/:id/delete/phone", api.UserDeletePhone(db))
	v1uEdit.POST("/:id/add/phone", api.UserAddPhone(db))
	v1uEdit.POST("/:id/add/address", api.UserAddAddress(db))
	v1uEdit.DELETE("/:id/delete/address", api.UserDeleteAddress(db))

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
