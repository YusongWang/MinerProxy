package global

import (
	"fmt"
	"math/big"
	pack "miner_proxy/pack/eth"
	"miner_proxy/utils"
	"sync"
	"time"

	"go.uber.org/zap"
)

type OnlineWorkers struct {
	Workers map[string]*Worker
	sync.Mutex
}

var GonlineWorkers = new(OnlineWorkers)

func init() {
	GonlineWorkers.Lock()
	GonlineWorkers.Workers = make(map[string]*Worker)
	GonlineWorkers.Unlock()
}

const (
	NOT_MINER    = 0
	MINER_LOGIN  = 1
	MINER_LOGOUT = 2
)

type Worker struct {
	Id            string               `json:"id"`
	Worker_name   string               `json:"worker_name"`
	Wallet        string               `json:"wallet"`
	Worker_idx    uint64               `json:"worker_idx"`
	Worker_share  uint64               `json:"worker_share"`
	Worker_reject uint64               `json:"worker_reject"`
	Report_hash   *big.Int             `json:"report_hash"`
	Login_time    time.Time            `json:"login_time"`
	Worker_diff   *big.Int             `json:"worker_diff"`
	Dev_idx       uint64               `json:"dev_idx"`
	Dev_diff      *big.Int             `json:"dev_diff"`
	Fee_idx       uint64               `json:"fee_idx"`
	Fee_diff      *big.Int             `json:"fee_diff"`
	Online        int                  `json:"online"`
	Ip            string               `json:"ip"`
	Ping          int                  `json:"ping"`
	Protocol      pack.StratumProtocol `json:"protocol"`
	AuthorizeStat pack.AuthorizeStat   `json:"authorize_stat"`
	OnlineTime    int64                `json:"online_time"`
}

func NewWorker(worker string, wallet string, id string, ip string) *Worker {
	return &Worker{
		Id:            id,
		Ip:            ip,
		Worker_name:   worker,
		Wallet:        wallet,
		Worker_idx:    0,
		Worker_share:  0,
		Worker_reject: 0,
		Report_hash:   new(big.Int).SetInt64(0),
		Login_time:    time.Now(),
		Worker_diff:   new(big.Int).SetInt64(0),
		Dev_idx:       0,
		Dev_diff:      new(big.Int).SetInt64(0),
		Fee_idx:       0,
		Fee_diff:      new(big.Int).SetInt64(0),
		Online:        NOT_MINER,
	}
}

func (w *Worker) FeeAdd() {
	w.Fee_idx++
}

func (w *Worker) SetFeeDiff(diff *big.Int) {
	w.Fee_diff = diff
}

func (w *Worker) GetFeeDiff() *big.Int {
	return w.Fee_diff
}

func (w *Worker) SetDevDiff(diff *big.Int) {
	w.Dev_diff = diff
}

func (w *Worker) GetDevDiff() *big.Int {
	return w.Dev_diff
}

func (w *Worker) SetDiff(diff *big.Int) {
	w.Worker_diff = diff
}

func (w *Worker) GetDiff() *big.Int {
	return w.Worker_diff
}

func (w *Worker) DevAdd() {
	w.Dev_idx++
}

func (w *Worker) AddShare() {
	w.Worker_share++
	//utils.Logger.Info("旷工下线.", zap.String("UUID", w.Id), zap.String("Worker", w.Worker_name), zap.String("Wallet", w.Wallet), zap.String("在线时长", w.Login_time.String()))
	utils.Logger.Info(fmt.Sprintf("Share #%d", w.Worker_share), zap.String("UUID", w.Id), zap.String("Worker", w.Worker_name), zap.String("Wallet", w.Wallet), zap.String("在线时长", w.Login_time.String()))
	//w.OnlineTime = humanize.Time(w.Login_time)
}

func (w *Worker) AddReject() {
	w.Worker_reject++
}

func (w *Worker) AddIndex() {
	w.Worker_idx++
}

func (w *Worker) GetIndex() uint64 {
	return w.Worker_idx
}

func (w *Worker) SetReportHash(hash *big.Int) {
	w.Report_hash = hash
}

func (w *Worker) SetPing(ping int) {
	w.Ping = ping
}

func (w *Worker) Logind(worker, wallet string) {
	w.Wallet = wallet
	w.Worker_name = worker
	w.Online = MINER_LOGIN
	utils.Logger.Info("登陆矿工.", zap.String("UUID", w.Id), zap.String("Worker", worker), zap.String("Wallet", wallet))
}

func (w *Worker) Logout() {
	w.Online = MINER_LOGOUT
	//utils.Logger.Info("旷工下线.", zap.String("UUID", w.Id), zap.String("Worker", w.Worker_name), zap.String("Wallet", w.Wallet), zap.String("在线时长", w.Login_time.String()))
}

func (w *Worker) IsOnline() bool {
	return w.Online == MINER_LOGIN
}

func (w *Worker) IsOffline() bool {
	return w.Online == MINER_LOGOUT
}
