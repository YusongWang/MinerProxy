package routes

import (
	"miner_proxy/asset"
	"miner_proxy/web/controllers"
	"miner_proxy/web/middleware/jwt"
	"time"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RegisterApiRouter(router *gin.Engine) {
	// models.InsertTest()
	// models.ReadMiners()
	// fs := assetfs.AssetFS{
	// 	Asset:     Asset,
	// 	AssetDir:  AssetDir,
	// 	AssetInfo: AssetInfo,
	// 	Prefix:    "",
	// 	Fallback:  "index.html",
	// }
	// router.StaticFS("/", &fs)
	fsStatic := assetfs.AssetFS{Asset: asset.Asset, AssetDir: asset.AssetDir, AssetInfo: asset.AssetInfo, Prefix: "dist/static", Fallback: "index.html"}
	fsThemes := assetfs.AssetFS{Asset: asset.Asset, AssetDir: asset.AssetDir, AssetInfo: asset.AssetInfo, Prefix: "dist/themes", Fallback: "index.html"}
	// fsCss := assetfs.AssetFS{Asset: asset.Asset, AssetDir: asset.AssetDir, AssetInfo: asset.AssetInfo, Prefix: "dist/css", Fallback: "index.html"}
	// fsFonts := assetfs.AssetFS{Asset: asset.Asset, AssetDir: asset.AssetDir, AssetInfo: asset.AssetInfo, Prefix: "dist/fonts", Fallback: "index.html"}
	// fsImg := assetfs.AssetFS{Asset: asset.Asset, AssetDir: asset.AssetDir, AssetInfo: asset.AssetInfo, Prefix: "dist/img", Fallback: "index.html"}
	// fsJs := assetfs.AssetFS{Asset: asset.Asset, AssetDir: asset.AssetDir, AssetInfo: asset.AssetInfo, Prefix: "dist/js", Fallback: "index.html"}
	fs := assetfs.AssetFS{Asset: asset.Asset, AssetDir: asset.AssetDir, AssetInfo: asset.AssetInfo, Prefix: "dist", Fallback: "index.html"}
	router.StaticFS("/static/", &fsStatic)
	router.StaticFS("/themes/", &fsThemes)
	// router.StaticFS("/css", &fsCss)
	// router.StaticFS("/fonts", &fsFonts)
	// router.StaticFS("/img", &fsImg)
	// router.StaticFS("/js", &fsJs)
	router.StaticFS("/favicon.ico", &fs)
	router.GET("/", func(c *gin.Context) {
		c.Writer.WriteHeader(200)
		indexHtml, _ := asset.Asset("dist/index.html")
		_, _ = c.Writer.Write(indexHtml)
		c.Writer.Header().Add("Accept", "text/html")
		c.Writer.Flush()
	})

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
		apiRouter.GET("/version", controllers.Version)
		apiRouter.GET("/announcement", controllers.Announcement)
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
