package controllers

import (
	"fmt"
	"math/big"
	"miner_proxy/global"
	"miner_proxy/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"go.uber.org/zap"
)

type Dashboard struct {
	PoolLength    int      `json:"pool_length"`
	OnlineWorker  int      `json:"online_worker"`
	OfflineWorker int      `json:"offline_worker"`
	TotalHash     *big.Int `json:"total_hash"`
	OnlineTime    string   `json:"online_time"`
	TotalShare    int64    `json:"total_shares"`
	TotalDiff     *big.Int `json:"total_diff"`
	FeeShares     int64    `json:"fee_shares"`
	FeeDiff       *big.Int `json:"fee_diff"`
	DevShares     int64    `json:"dev_shares"`
	DevDiff       *big.Int `json:"dev_diff"`
}

func newDashborad() *Dashboard {

	return &Dashboard{
		PoolLength:    0,
		OnlineWorker:  0,
		OfflineWorker: 0,
		TotalHash:     new(big.Int).SetInt64(0),
		OnlineTime:    "",
		TotalShare:    0,
		TotalDiff:     new(big.Int).SetInt64(0),
		FeeShares:     0,
		FeeDiff:       new(big.Int).SetInt64(0),
		DevDiff:       new(big.Int).SetInt64(0),
	}
}

// 首页展示数据接口
func Home(c *gin.Context) {

	eth := newDashborad()
	etc := newDashborad()

	eth_res := make(map[string]interface{})
	etc_res := make(map[string]interface{})

	for _, app := range global.ManageApp.Config {
		if app.Coin == "ETH" {
			eth.PoolLength++
			for _, w := range global.OnlinePools[app.ID] {
				if w.IsOnline() {
					eth.OnlineWorker++
					eth.TotalShare = eth.TotalShare + int64(w.Worker_share)
					eth.FeeShares = eth.FeeShares + int64(w.Fee_idx)
					eth.DevShares = eth.DevShares + int64(w.Dev_idx)
					eth.TotalHash = new(big.Int).Add(eth.TotalHash, w.Report_hash)
					eth.TotalDiff = new(big.Int).Add(eth.TotalDiff, w.Worker_diff)
					eth.FeeDiff = new(big.Int).Add(eth.FeeDiff, w.Fee_diff)
					eth.DevDiff = new(big.Int).Add(eth.DevDiff, w.Dev_diff)

				} else if w.IsOffline() {
					eth.OfflineWorker++
				}
			}
		}

		if app.Coin == "ETC" {
			etc.PoolLength++
			for _, w := range global.OnlinePools[app.ID] {

				if w.IsOnline() {
					etc.OnlineWorker++
					etc.TotalShare = etc.TotalShare + int64(w.Worker_share)
					etc.FeeShares = etc.FeeShares + int64(w.Fee_idx)
					etc.DevShares = etc.DevShares + int64(w.Dev_idx)
					etc.TotalHash = new(big.Int).Add(etc.TotalHash, w.Report_hash)

					etc.TotalDiff = new(big.Int).Add(etc.TotalDiff, w.Worker_diff)
					etc.FeeDiff = new(big.Int).Add(etc.FeeDiff, w.Fee_diff)
					etc.DevDiff = new(big.Int).Add(etc.DevDiff, w.Dev_diff)

				} else if w.IsOffline() {
					etc.OfflineWorker++
				}
			}
		}
	}

	eth_res["online_worker"] = eth.OnlineWorker
	eth_res["pool_length"] = eth.PoolLength
	eth_res["offline_worker"] = eth.OfflineWorker
	eth_res["total_hash"] = eth.TotalHash
	eth_res["online_time"] = global.Start_Time
	eth_res["total_shares"] = eth.TotalShare
	if eth.OnlineWorker > 0 {
		eth_res["total_diff"] = new(big.Int).Div(eth.TotalDiff, new(big.Int).SetInt64(int64(eth.OnlineWorker)))
		eth_res["fee_diff"] = new(big.Int).Div(eth.FeeDiff, new(big.Int).SetInt64(int64(eth.OnlineWorker)))
		eth_res["dev_diff"] = new(big.Int).Div(eth.DevDiff, new(big.Int).SetInt64(int64(eth.OnlineWorker)))
	} else {
		eth_res["total_diff"] = "0"
		eth_res["fee_diff"] = "0"
		eth_res["dev_diff"] = "0"
	}

	eth_res["fee_shares"] = eth.FeeShares
	eth_res["fee_rate"] = fmt.Sprintf("%.2f", float64(eth.FeeShares)/float64(eth.TotalShare)*100.0)
	eth_res["dev_shares"] = eth.DevShares
	eth_res["dev_rate"] = fmt.Sprintf("%.2f", float64(eth.DevShares)/float64(eth.TotalShare)*100.0)

	etc_res["online_worker"] = etc.OnlineWorker
	etc_res["pool_length"] = etc.PoolLength
	etc_res["offline_worker"] = etc.OfflineWorker
	etc_res["total_hash"] = etc.TotalHash
	etc_res["online_time"] = global.Start_Time
	etc_res["total_shares"] = etc.TotalShare

	etc_res["fee_shares"] = etc.FeeShares

	etc_res["fee_rate"] = fmt.Sprintf("%.2f", float64(etc.FeeShares)/float64(etc.TotalShare)*100.0)
	etc_res["dev_shares"] = etc.DevShares

	etc_res["dev_rate"] = fmt.Sprintf("%.2f", float64(etc.DevShares)/float64(etc.TotalShare)*100.0)
	if etc.OnlineWorker > 0 {
		etc_res["dev_diff"] = new(big.Int).Div(etc.DevDiff, new(big.Int).SetInt64(int64(etc.OnlineWorker)))
		etc_res["fee_diff"] = new(big.Int).Div(etc.FeeDiff, new(big.Int).SetInt64(int64(etc.OnlineWorker)))
		etc_res["total_diff"] = new(big.Int).Div(etc.TotalDiff, new(big.Int).SetInt64(int64(etc.OnlineWorker)))
	} else {
		etc_res["dev_diff"] = "0"
		etc_res["fee_diff"] = "0"
		etc_res["total_diff"] = "0"
	}

	var data = map[string]map[string]interface{}{"ETH": eth_res, "ETC": etc_res}

	c.JSON(200, gin.H{
		"data": data,
		"msg":  "",
		"code": 200,
	})
}

// 系统性能接口
func System(c *gin.Context) {
	utils.Logger.Info("cpu", zap.Any("Cpu", GetCpuPercent()))
	utils.Logger.Info("mem", zap.Any("mem", GetMemPercent()))
	c.JSON(200, gin.H{
		"msg": "pong",
	})
}

func GetCpuPercent() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	return percent[0]
}

func GetMemPercent() float64 {
	memInfo, _ := mem.VirtualMemory()
	return memInfo.UsedPercent
}
