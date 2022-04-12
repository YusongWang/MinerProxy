package pack

import "math/big"

type Job struct {
	Job [][]string
}

type Worker struct {
	Worker_name   string
	Wallet        string
	Worker_idx    uint64
	Worker_share  uint64
	Worker_reject uint64
	Worker_diff   *big.Int
	Dev_idx       uint64
	Dev_diff      *big.Int
	Fee_idx       uint64
	Fee_diff      *big.Int
	IsOnline      bool
}

func NewWorker(worker string, wallet string) *Worker {

	return &Worker{
		Worker_name:   worker,
		Wallet:        wallet,
		Worker_idx:    0,
		Worker_share:  0,
		Worker_reject: 0,
		Worker_diff:   new(big.Int).SetInt64(0),
		Dev_idx:       0,
		Dev_diff:      new(big.Int).SetInt64(0),
		Fee_idx:       0,
		Fee_diff:      new(big.Int).SetInt64(0),
		IsOnline:      false,
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

func (w *Worker) Logind(worker, wallet string) {
	w.Wallet = wallet
	w.Worker_name = worker
	w.IsOnline = true
}

func (w *Worker) Logout() {
	w.IsOnline = false
}
