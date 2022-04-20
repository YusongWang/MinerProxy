package controllers

import (
	"miner_proxy/global"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 单个矿池的矿工列表
func MinerList(c *gin.Context) {
	id_str := c.PostForm("id")

	id, err := strconv.Atoi(id_str)
	if err != nil {
		c.JSON(200, gin.H{
			"msg":  "矿池ID未选择",
			"code": 300,
		})
		return
	}

	c.JSON(200, gin.H{
		"data": global.OnlinePools[id],
		"msg":  "",
		"code": 200,
	})
}

// 单个矿工详情
func MinerDetail(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
