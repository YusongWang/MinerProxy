package controllers

import (
	"fmt"
	"miner_proxy/global"

	"github.com/gin-gonic/gin"
)

type WorkerList struct {
	Name      string
	TotalHash string
	Online    int
	Offline   int    `json:"off_line"`
	Coin      string `json:"coin"`
	Port      int
	Protocol  string `json:"protocol"`
	IsRun     bool   `json:"is_run"`
}

// 展示矿池列表 在线和不在线的
func PoolList(c *gin.Context) {
	var list []WorkerList
	for _, l := range global.ManageApp.Config {
		temp := WorkerList{
			Name:      l.Worker,
			TotalHash: "",
			Online:    0,
			Offline:   0,
			Coin:      l.Coin,
			Port:      l.TLS,
			Protocol:  "SSL",
			IsRun:     l.Online,
		}

		if global.OnlinePools[l.ID] != nil {
			temp.Online = len(global.OnlinePools[l.ID])
			for _, worker := range global.OnlinePools[l.ID] {
				//TODO HASHRATE 未记录。
				fmt.Println(worker)
			}
		}

		list = append(list, temp)
	}

	c.JSON(200, gin.H{
		"data":    list,
		"message": "",
		"code":    200,
	})
}
