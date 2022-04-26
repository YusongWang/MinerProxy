package controllers

import (
	"encoding/json"
	"miner_proxy/global"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 展示矿池列表 在线和不在线的
func SetPass(c *gin.Context) {
	old := c.PostForm("oldpass")
	new := c.PostForm("pass")
	if new == "" {
		c.JSON(200, gin.H{
			"data":    "",
			"message": "请传入密码",
			"code":    301,
		})
		return
	}

	if old == global.ManageApp.Web.Password {
		global.ManageApp.Web.Password = new

		config_json, err := json.Marshal(global.ManageApp)
		if err != nil {
			c.JSON(200, gin.H{
				"data":    "",
				"message": "格式化配置文件失败",
				"code":    301,
			})
			return
		}
		config_file, err := os.OpenFile("config.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0777)
		if err != nil {
			c.JSON(200, gin.H{
				"data":    "",
				"message": "打开配置文件失败",
				"code":    301,
			})
			return
		}

		config_file.Write(config_json)
		config_file.Close()
		c.JSON(200, gin.H{
			"data":    "",
			"message": "修改成功，请刷新网页.",
			"code":    200,
		})
	} else {
		c.JSON(200, gin.H{
			"data":    "",
			"message": "旧密码不正确",
			"code":    301,
		})
	}
}

// 展示矿池列表 在线和不在线的
func SetPort(c *gin.Context) {
	new := c.PostForm("port")
	if new == "" {
		c.JSON(200, gin.H{
			"data":    "",
			"message": "请传入端口号",
			"code":    301,
		})
		return
	}

	newport, err := strconv.Atoi(new)
	if err != nil {
		c.JSON(200, gin.H{
			"data":    "",
			"message": "端口号不正确",
			"code":    301,
		})
		return
	}
	if newport <= 1000 || newport >= 65535 {
		c.JSON(200, gin.H{
			"data":    "",
			"message": "端口号不正确",
			"code":    301,
		})
		return
	}

	if newport != global.ManageApp.Web.Port {
		global.ManageApp.Web.Port = newport

		config_json, err := json.Marshal(global.ManageApp)
		if err != nil {
			c.JSON(200, gin.H{
				"data":    "",
				"message": "格式化配置文件失败",
				"code":    301,
			})
			return
		}
		config_file, err := os.OpenFile("config.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0777)
		if err != nil {
			c.JSON(200, gin.H{
				"data":    "",
				"message": "打开配置文件失败",
				"code":    301,
			})
			return
		}

		config_file.Write(config_json)
		config_file.Close()
		c.JSON(200, gin.H{
			"data":    "",
			"message": "修改成功，请刷新网页.",
			"code":    200,
		})
	} else {
		c.JSON(200, gin.H{
			"data":    "",
			"message": "修改成功",
			"code":    200,
		})
	}

}
