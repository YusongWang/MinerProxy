package controllers

import (
	"encoding/json"
	"math/big"
	"miner_proxy/global"
	"miner_proxy/utils"
	"os"

	"github.com/gin-gonic/gin"
)

type WorkerList struct {
	Name      string   `json:"name"`
	ID        int      `json:"id"`
	TotalHash *big.Int `json:"total_hash"`
	Online    int      `json:"online"`
	Offline   int      `json:"off_line"`
	Coin      string   `json:"coin"`
	Port      int      `json:"port"`
	Protocol  string   `json:"protocol"`
	IsRun     bool     `json:"is_run"`
}

// 展示矿池列表 在线和不在线的
func PoolList(c *gin.Context) {
	var list []WorkerList
	for _, l := range global.ManageApp.Config {
		temp := WorkerList{
			Name:      l.Worker,
			TotalHash: new(big.Int).SetInt64(0),
			ID:        l.ID,
			Online:    0,
			Offline:   0,
			Coin:      l.Coin,
			Port:      l.TLS,
			Protocol:  "SSL",
			IsRun:     l.Online,
		}

		if global.OnlinePools[l.ID] != nil {
			for _, worker := range global.OnlinePools[l.ID] {
				if worker.IsOnline() {
					temp.Online++
					temp.TotalHash = new(big.Int).Add(temp.TotalHash, worker.Report_hash)
				} else {
					temp.Offline++
				}
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

func CreatePool(c *gin.Context) {
	var config utils.Config
	err := c.BindJSON(&config)
	if err != nil {
		c.JSON(200, gin.H{
			"data":    "",
			"message": "解析参数失败" + err.Error(),
			"code":    301,
		})
		return
	}

	global.ManageApp.Config = append(global.ManageApp.Config, config)

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
		"message": "添加成功",
		"code":    200,
	})
}
