package controllers

import "github.com/gin-gonic/gin"

// 首页展示数据接口
func Home(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

//系统性能接口
func System(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
