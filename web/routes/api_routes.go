package routes

import (
	//"morningo/filters/auth"

	"miner_proxy/web/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterApiRouter(router *gin.Engine) {
	apiRouter := router.Group("api")
	{
		apiRouter.GET("/dashbroad", controllers.Home)
		apiRouter.GET("/system", controllers.System)
		apiRouter.GET("/miner/detail", controllers.MinerDetail)
		apiRouter.GET("/miners", controllers.MinerList)
		apiRouter.GET("/pools", controllers.PoolList)
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
