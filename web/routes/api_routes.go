package routes

import (
	"miner_proxy/web/controllers"
	"miner_proxy/web/middleware/jwt"
	"miner_proxy/web/models"

	"github.com/gin-gonic/gin"
)

func RegisterApiRouter(router *gin.Engine) {

	router.Use(gin.Logger())

	models.InsertTest()
	models.ReadMiners()
	router.POST("auth/login", controllers.Login)

	apiRouter := router.Group("api")
	apiRouter.Use(jwt.JWT())
	{
		apiRouter.POST("/dashborad", controllers.Home)
		apiRouter.POST("/system", controllers.System)
		apiRouter.POST("/miner/detail", controllers.MinerDetail)
		apiRouter.POST("/miners", controllers.MinerList)
		apiRouter.POST("/pools", controllers.PoolList)
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
