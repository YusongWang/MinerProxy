package controllers

import (
	"miner_proxy/global"
	"miner_proxy/utils"

	"github.com/gin-gonic/gin"
)

type Password struct {
	Password string `json:"password"`
}

// 登录接口,判断密码是否与config.json一致。不一致则进制登录。登录后写入jsonwebtoken . 用中间件进行判断。
func Login(c *gin.Context) {
	data := make(map[string]interface{})

	var p Password
	err := c.BindJSON(&p)
	if err != nil {
		c.JSON(200, gin.H{
			"data": "",
			"msg":  "解析参数失败" + err.Error(),
			"code": 301,
		})
		return
	}

	if p.Password == "" {
		c.JSON(200, gin.H{
			"code": 301,
			"msg":  "请输入密码",
			"data": data,
		})
		return
	}

	if p.Password != global.ManageApp.Web.Password {
		c.JSON(200, gin.H{
			"code": 302,
			"msg":  "密码不正确",
			"data": data,
		})
		return
	}

	token, err := utils.GenerateToken(p.Password)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 303,
			"msg":  "秘钥生成错误",
			"data": data,
		})
		return
	} else {
		data["token"] = token
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "",
		"data": data,
	})
}

// 设置登录密码和端口。然后保存到config.json
func WebConfig(c *gin.Context) {

	c.JSON(200, gin.H{
		"msg": "pong",
	})
}
