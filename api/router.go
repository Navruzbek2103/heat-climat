package api

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	swaggerfile "github.com/swaggo/files"
	_ "gitlab.com/climate.uz/api/docs"
	v1 "gitlab.com/climate.uz/api/handler/v1"
	"gitlab.com/climate.uz/api/middleware"
	"gitlab.com/climate.uz/api/tokens"
	"gitlab.com/climate.uz/config"
	"gitlab.com/climate.uz/internal/controller/storage"
	"gitlab.com/climate.uz/pkg/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Options struct {
	Conf           config.Config
	Logger         logger.Logger
	Storage        storage.StorageI
	CasbinEnforcer *casbin.Enforcer
}

// New ...
// @title           Climate_uz api
// @version         1.0
// @description     This is climate_uz server api server
// @termsOfService  2 term climate_uz

// @contact.name   	Murtazoxon
// @contact.url    	https://t.me/murtazokhon_gofurov
// @contact.email  	gofurovmurtazoxon@gmail.com

// @host      		100.26.188.92:7070
// @BasePath  		/v1

// @securityDefinitions.apikey BearerAuth
// @in 		header
// @name 	Authorization
func New(option Options) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	corConfig := cors.DefaultConfig()
	corConfig.AllowAllOrigins = true
	corConfig.AllowCredentials = true
	corConfig.AllowHeaders = []string{"*"}
	corConfig.AllowBrowserExtensions = true
	corConfig.AllowMethods = []string{"*"}
	router.Use(cors.New(corConfig))

	jwtHandler := tokens.JWTHandler{
		SigninKey: option.Conf.SigninKey,
		Log:       option.Logger,
	}
	handlerV1 := v1.New(&v1.HandlerV1Config{
		Logger:     option.Logger,
		Storage:    option.Storage,
		Cfg:        option.Conf,
		JwtHandler: jwtHandler,
	})

	router.Use(middleware.NewAuth(option.CasbinEnforcer, jwtHandler, config.Load()))

	router.MaxMultipartMemory = 8 << 20 // 8 Mib

	router.Static("/media", "./media")

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "App is running",
		})
	})
	api := router.Group("/v1")

	// Brand api
	api.POST("/brand", handlerV1.CreateBrand)
	api.PATCH("/brand", handlerV1.UpdateBrandInfo)
	api.GET("/brand/:id", handlerV1.GetBrandId)
	api.POST("/brand-category", handlerV1.CreateBrandCategory)

	// Category apis
	api.POST("/category", handlerV1.CreateCategory)
	api.GET("/category/:id", handlerV1.GetCategory)
	api.PATCH("/category", handlerV1.EditCategoryName)
	api.GET("/categories", handlerV1.GetAllGategory)
	api.DELETE("category/:id", handlerV1.DeleteCategory)

	// User apis
	api.POST("/user", handlerV1.CreateUser)
	api.GET("/user/:id", handlerV1.GetUser)
	api.GET("/users", handlerV1.GetAllUsers)
	api.PATCH("/user", handlerV1.UpdateUser)
	api.DELETE("/user/:id", handlerV1.DeleteUser)

	// Product apis
	api.POST("/product", handlerV1.CreateProduct)
	api.GET("/product/:id", handlerV1.GetProductInfo)
	api.PUT("/product", handlerV1.UpdateProduct)
	api.GET("/products", handlerV1.GetAllProduct)
	api.DELETE("/product/:id", handlerV1.DeleteProductInfo)

	// News product apis
	api.POST("/new-product", handlerV1.CreateNewProduct)
	api.GET("/new-product/:id", handlerV1.GetNewProduct)
	api.GET("/new-products", handlerV1.GetAllNewsProducts)
	api.PATCH("/new-product", handlerV1.UpdateNewsProduct)
	api.DELETE("/new-product/:id", handlerV1.DeleteNewsProduct)

	// admin api
	api.POST("/admin", handlerV1.AddAdmin)
	api.GET("/admin/:username/:password", handlerV1.AdminLogin)
	api.GET("/admin/profile", handlerV1.GetAdminProfile)
	api.GET("/admins", handlerV1.GetAllAdmin)
	api.PATCH("/admin", handlerV1.UpdateAdminInfo)
	api.DELETE("/admin/:id", handlerV1.DeleteAdminInfo)

	// upload image
	api.POST("/image-upload", handlerV1.UploadMedia)
	url := ginSwagger.URL("swagger/doc.json")
	api.GET("swagger/*any", ginSwagger.WrapHandler(swaggerfile.Handler, url))

	return router
}
