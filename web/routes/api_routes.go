package routes

import (
	"miner_proxy/web/controllers"
	"miner_proxy/web/middleware/jwt"
	"miner_proxy/web/models"
	"time"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterApiRouter(router *gin.Engine) {
	models.InsertTest()
	models.ReadMiners()
	fs := assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
		Prefix:    "",
		Fallback:  "index.html",
	}
	router.StaticFS("/", &fs)
	router.Use(gin.Logger())

	//	router.StaticFS("/", http.FileServer(
	//		&assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: "dist"}))
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	router.POST("auth/login", controllers.Login)

	apiRouter := router.Group("api")
	apiRouter.Use(jwt.JWT())
	{
		apiRouter.GET("/dashborad", controllers.Home)
		//apiRouter.POST("/system", controllers.System)
		//apiRouter.POST("/miner/detail", controllers.MinerDetail)
		apiRouter.GET("/miners/:id", controllers.MinerList)

		apiRouter.POST("/setpass", controllers.SetPass)
		apiRouter.POST("/setport", controllers.SetPort)

		apiRouter.GET("/pool", controllers.PoolList)
		apiRouter.POST("/pool", controllers.CreatePool)
		apiRouter.GET("/pool/:id", controllers.GetPool)
		apiRouter.POST("/pool/:id", controllers.UpdatePool)
		apiRouter.DELETE("/pool/:id", controllers.DeletePool)
	}

	// api := router.Group("/api")
	// api.GET("/index", controllers.IndexApi)
	// api.GET("/cookie/set/:userid", controllers.CookieSetExample)

	// cookie auth middleware
	// api.Use(auth.Middleware(auth.CookieAuthDriverKey))
	// {
	// 	api.GET("/orm", controllers.OrmExample)
	// 	api.GET("/store", controllers.StoreExample)
	// 	api.GET("/db", controllers.DBExample)
	// 	api.GET("/cookie/get", controllers.CookieGetExample)
	// }

	// jwtApi := router.Group("/api")
	// jwtApi.GET("/jwt/set/:userid", controllers.JwtSetExample)

	// // jwt auth middleware
	// jwtApi.Use(auth.Middleware(auth.JwtAuthDriverKey))
	// {
	// 	jwtApi.GET("/jwt/get", controllers.JwtGetExample)
	// }
}
