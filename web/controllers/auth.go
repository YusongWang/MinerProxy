package controllers

import "github.com/gin-gonic/gin"

// 登录接口,判断密码是否与config.json一致。不一致则进制登录。登录后写入jsonwebtoken . 用中间件进行判断。
func Login(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

// 设置登录密码和端口。然后保存到config.json
func WebConfig(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
