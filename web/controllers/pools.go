package controllers

import "github.com/gin-gonic/gin"

// 展示矿池列表 在线和不在线的
func PoolList(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
