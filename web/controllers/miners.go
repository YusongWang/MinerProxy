package controllers

import (
	"fmt"
	"miner_proxy/global"
	"miner_proxy/pack"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 单个矿池的矿工列表
func MinerList(c *gin.Context) {
	id_str := c.Param("id")

	id, err := strconv.Atoi(id_str)
	if err != nil {
		c.JSON(200, gin.H{
			"msg":  "矿池ID未选择",
			"code": 300,
		})
		return
	}

	var res []pack.Worker
	if len(global.OnlinePools[id]) > 0 {
		for _, miner := range global.OnlinePools[id] {
			if miner.IsOnline() {
				miner.OnlineTime = fmt.Sprintf("%d", miner.Login_time.Unix())
				res = append(res, miner)
			}
		}
	}

	c.JSON(200, gin.H{
		"data": res,
		"msg":  "",
		"code": 200,
	})
}

// 单个矿工详情
func MinerDetail(c *gin.Context) {
	c.JSON(200, gin.H{
		"msg": "pong",
	})
}
