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
	Name          string   `json:"name"`
	ID            int      `json:"id"`
	TotalHash     *big.Int `json:"total_hash"`
	Online        int      `json:"online"`
	Offline       int      `json:"off_line"`
	Coin          string   `json:"coin"`
	Port          int      `json:"port"`
	Protocol      string   `json:"protocol"`
	IsRun         bool     `json:"is_run"`
	OnlineWorker  int      `json:"online_worker"`
	OfflineWorker int      `json:"offline_worker"`
	OnlineTime    string   `json:"online_time"`
	TotalShare    int64    `json:"total_shares"`
	TotalDiff     *big.Int `json:"total_diff"`
	FeeShares     int64    `json:"fee_shares"`
	FeeDiff       *big.Int `json:"fee_diff"`
	DevShares     int64    `json:"dev_shares"`
	DevDiff       *big.Int `json:"dev_diff"`
}

// 展示矿池列表 在线和不在线的
func PoolList(c *gin.Context) {
	var list []WorkerList
	for _, l := range global.ManageApp.Config {
		temp := WorkerList{
			Name:          l.Worker,
			TotalHash:     new(big.Int).SetInt64(0),
			ID:            l.ID,
			Online:        0,
			Offline:       0,
			Coin:          l.Coin,
			Port:          l.TLS,
			Protocol:      "SSL",
			IsRun:         l.Online,
			OnlineWorker:  0,
			OfflineWorker: 0,
			OnlineTime:    "",
			TotalShare:    0,
			TotalDiff:     new(big.Int).SetInt64(0),
			FeeShares:     0,
			FeeDiff:       new(big.Int).SetInt64(0),
			DevDiff:       new(big.Int).SetInt64(0),
		}

		if global.OnlinePools[l.ID] != nil {
			for _, w := range global.OnlinePools[l.ID] {
				if w.IsOnline() {
					temp.Online++
					temp.TotalHash = new(big.Int).Add(temp.TotalHash, w.Report_hash)
					temp.FeeShares = temp.FeeShares + int64(w.Fee_idx)
					temp.DevShares = temp.DevShares + int64(w.Dev_idx)
					temp.TotalHash = new(big.Int).Add(temp.TotalHash, w.Report_hash)
					temp.TotalDiff = new(big.Int).Add(temp.TotalDiff, w.Worker_diff)
					temp.FeeDiff = new(big.Int).Add(temp.FeeDiff, w.Fee_diff)
					temp.DevDiff = new(big.Int).Add(temp.DevDiff, w.Dev_diff)
				} else {
					temp.Offline++
				}
			}
		}

		if temp.OnlineWorker > 0 {
			temp.DevDiff = new(big.Int).Div(temp.DevDiff, new(big.Int).SetInt64(int64(temp.OnlineWorker)))
			temp.FeeDiff = new(big.Int).Div(temp.FeeDiff, new(big.Int).SetInt64(int64(temp.OnlineWorker)))
			temp.TotalDiff = new(big.Int).Div(temp.TotalDiff, new(big.Int).SetInt64(int64(temp.OnlineWorker)))
		}

		list = append(list, temp)
	}

	c.JSON(200, gin.H{
		"data": list,
		"msg":  "",
		"code": 200,
	})
}

func CreatePool(c *gin.Context) {
	var config utils.Config
	err := c.BindJSON(&config)
	if err != nil {
		c.JSON(200, gin.H{
			"data": "",
			"msg":  "解析参数失败" + err.Error(),
			"code": 301,
		})
		return
	}

	global.ManageApp.Config = append(global.ManageApp.Config, config)

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
		"msg":  "添加成功",
		"code": 200,
	})
}
