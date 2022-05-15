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
	//fs := assetfs.AssetFS{Asset: asset.Asset, AssetDir: asset.AssetDir, AssetInfo: asset.AssetInfo, Prefix: "dist", Fallback: "index.html"}
	router.StaticFS("/static/", &fsStatic)
	router.StaticFS("/themes/", &fsThemes)
	// router.StaticFS("/css", &fsCss)
	// router.StaticFS("/fonts", &fsFonts)
	// router.StaticFS("/img", &fsImg)
	// router.StaticFS("/js", &fsJs)
	router.StaticFS("/favicon.ico", &fs)
	router.StaticFS("/icon.png", &fs)

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
	router.GET("logger", controllers.Logger)

	apiRouter := router.Group("api")
	apiRouter.Use(jwt.JWT())
	{
		apiRouter.GET("/version", controllers.Version)
		apiRouter.GET("/announcement", controllers.Announcement)
		apiRouter.GET("/dashborad", controllers.Home)

		apiRouter.GET("/system_chart", controllers.SystemChart)
		apiRouter.GET("/worker_chart/:coin", controllers.Worker)

		apiRouter.GET("/miners/:id", controllers.MinerList)

		apiRouter.POST("/setpass", controllers.SetPass)
		apiRouter.POST("/setport", controllers.SetPort)

		apiRouter.GET("/pool", controllers.PoolList)
		apiRouter.POST("/pool", controllers.CreatePool)
		apiRouter.GET("/pool/:id", controllers.GetPool)
		apiRouter.POST("/pool/:id", controllers.UpdatePool)
		apiRouter.DELETE("/pool/:id", controllers.DeletePool)
	}
}
