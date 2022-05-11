package controllers

import (
	"miner_proxy/global"

	"miner_proxy/web/logics"
	"miner_proxy/web/models"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/levigross/grequests"
)

// 首页展示数据接口
func Home(c *gin.Context) {
	eth_res, etc_res := logics.ClacDashborad()
	var data = map[string]map[string]interface{}{"ETH": eth_res, "ETC": etc_res}

	c.JSON(200, gin.H{
		"data": data,
		"msg":  "",
		"code": 200,
	})
}

func Version(c *gin.Context) {
	//https://cdn.jsdelivr.net/gh/OutlawQAQ/MinerProxy@main/others/Version.md
	resp, err := grequests.Get("https://cdn.jsdelivr.net/gh/OutlawQAQ/MinerProxy@main/others/Version.md", nil)
	if err != nil {
		c.JSON(200, gin.H{
			"data": err.Error(),
			"msg":  "",
			"code": 300,
		})
	}

	doc := resp.String()
	doc = strings.Replace(doc, "{BUILD_VERSION}", global.Version, 1)

	c.JSON(200, gin.H{
		"data": doc,
		"msg":  "",
		"code": 200,
	})
}

func Announcement(c *gin.Context) {
	//https://cdn.jsdelivr.net/gh/OutlawQAQ/MinerProxy@main/others/announcement.md
	resp, err := grequests.Get("https://cdn.jsdelivr.net/gh/OutlawQAQ/MinerProxy@main/others/announcement.md", nil)
	if err != nil {
		c.JSON(200, gin.H{
			"data": err.Error(),
			"msg":  "",
			"code": 300,
		})
	}

	c.JSON(200, gin.H{
		"data": resp.String(),
		"msg":  "",
		"code": 200,
	})
}

func SystemChart(c *gin.Context) {
	resp, err := models.GetSys()
	if err != nil {
		c.JSON(200, gin.H{
			"data": err.Error(),
			"msg":  "fetch error",
			"code": 300,
		})
		return
	}

	c.JSON(200, gin.H{
		"data": resp,
		"msg":  "",
		"code": 200,
	})
}

func Worker(c *gin.Context) {
	coin := c.Param("coin")

	resp, err := models.GetWorker(coin)
	if err != nil {
		c.JSON(200, gin.H{
			"data": err.Error(),
			"msg":  "fetch error",
			"code": 300,
		})
		return
	}

	c.JSON(200, gin.H{
		"data": resp,
		"msg":  "",
		"code": 200,
	})
}
