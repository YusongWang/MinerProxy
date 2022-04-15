package controllers

import "github.com/gin-gonic/gin"

// 单个矿池的矿工列表
func MinerList(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// 单个矿工详情
func MinerDetail(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
