package controllers

import (
	"encoding/json"
	"miner_proxy/global"
	"os"

	"github.com/gin-gonic/gin"
)

type Passwrod struct {
	OldPass string `json:"oldpass"`
	Pass    string `json:"pass"`
}

// 展示矿池列表 在线和不在线的
func SetPass(c *gin.Context) {
	var pass Passwrod
	err := c.BindJSON(&pass)
	if err != nil {
		c.JSON(200, gin.H{
			"data": "",
			"msg":  "解析参数失败" + err.Error(),
			"code": 301,
		})
		return
	}

	if pass.OldPass == "" || pass.Pass == "" {
		c.JSON(200, gin.H{
			"data": "",
			"msg":  "请传入密码",
			"code": 301,
		})
		return
	}

	if pass.OldPass == global.ManageApp.Web.Password {
		global.ManageApp.Web.Password = pass.Pass
		config_json, err := json.Marshal(global.ManageApp)
		if err != nil {
			c.JSON(200, gin.H{
				"data": "",
				"msg":  "格式化配置文件失败",
				"code": 301,
			})
			return
		}
		config_file, err := os.OpenFile("config.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0777)
		if err != nil {
			c.JSON(200, gin.H{
				"data": "",
				"msg":  "打开配置文件失败",
				"code": 301,
			})
			return
		}

		config_file.Write(config_json)
		config_file.Close()
		c.JSON(200, gin.H{
			"data": "",
			"msg":  "修改成功，请刷新网页.",
			"code": 200,
		})
	} else {
		c.JSON(200, gin.H{
			"data": "",
			"msg":  "旧密码不正确",
			"code": 301,
		})
	}
}

type Port struct {
	Port int `json:"port"`
}

// 展示矿池列表 在线和不在线的
func SetPort(c *gin.Context) {
	var p Port
	err := c.BindJSON(&p)
	if err != nil {
		c.JSON(200, gin.H{
			"data": "",
			"msg":  "解析参数失败" + err.Error(),
			"code": 301,
		})
		return
	}

	if p.Port == 0 {
		c.JSON(200, gin.H{
			"data": "",
			"msg":  "请传入端口号",
			"code": 301,
		})
		return
	}

	if p.Port <= 0 || p.Port >= 65535 {
		c.JSON(200, gin.H{
			"data": "",
			"msg":  "端口号不正确",
			"code": 301,
		})
		return
	}

	if p.Port != global.ManageApp.Web.Port {
		global.ManageApp.Web.Port = p.Port

		config_json, err := json.Marshal(global.ManageApp)
		if err != nil {
			c.JSON(200, gin.H{
				"data": "",
				"msg":  "格式化配置文件失败",
				"code": 301,
			})
			return
		}
		config_file, err := os.OpenFile("config.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0777)
		if err != nil {
			c.JSON(200, gin.H{
				"data": "",
				"msg":  "打开配置文件失败",
				"code": 301,
			})
			return
		}

		config_file.Write(config_json)
		config_file.Close()
		c.JSON(200, gin.H{
			"data": "",
			"msg":  "修改成功，请刷新网页.",
			"code": 200,
		})
	} else {
		c.JSON(200, gin.H{
			"data": "",
			"msg":  "修改成功",
			"code": 200,
		})
	}

}
