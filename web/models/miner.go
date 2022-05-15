package models

import (
	"fmt"
	"math/big"
	"miner_proxy/global"

	memdb "github.com/hashicorp/go-memdb"
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

func NewDashborad() *Dashboard {

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

type WorkerChart struct {
	Time     int64    `json:"time"`
	Hashrate *big.Int `json:"hashrate"`
	Online   int      `json:"online"`
	Offline  int      `json:"offline"`
	Coin     string   `json:"coin"`
}

type SystemChart struct {
	Time int64   `json:"time"`
	Mem  float64 `json:"memory"`
	Cpu  float64 `json:"cpu"`
}

var Chart *memdb.MemDB

func init() {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"chart_ETC": &memdb.TableSchema{
				Name: "chart_ETC",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.IntFieldIndex{Field: "Time"},
					},
					"coin": &memdb.IndexSchema{
						Name:    "coin",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Coin"},
					},
				},
			},
			"chart_ETH": &memdb.TableSchema{
				Name: "chart_ETH",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.IntFieldIndex{Field: "Time"},
					},
					"coin": &memdb.IndexSchema{
						Name:    "coin",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Coin"},
					},
				},
			},
			"chart_system": &memdb.TableSchema{
				Name: "chart_system",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.IntFieldIndex{Field: "Time"},
					},
				},
			},
		},
	}

	// Create a new data base
	var err error
	Chart, err = memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}
}

type Miner struct {
	global.Worker
}

func NewMiner(w global.Worker) (*Miner, error) {
	return &Miner{w}, nil
}

func InsertWorkerETH(w WorkerChart) error {
	txn := Chart.Txn(true)

	if err := txn.Insert("chart_ETH", w); err != nil {
		return err
	}

	txn.Commit()
	return nil
}

func InsertWorkerETC(w WorkerChart) error {
	txn := Chart.Txn(true)

	if err := txn.Insert("chart_ETC", w); err != nil {
		return err
	}

	txn.Commit()
	return nil
}

func InsertSys(s SystemChart) error {
	txn := Chart.Txn(true)

	if err := txn.Insert("chart_system", s); err != nil {
		return err
	}

	txn.Commit()
	return nil
}

func GetSys() ([]SystemChart, error) {
	txn := Chart.Txn(false)
	defer txn.Abort()

	var system []SystemChart

	it, err := txn.Get("chart_system", "id")
	if err != nil {
		return nil, err
	}

	for obj := it.Next(); obj != nil; obj = it.Next() {
		p := obj.(SystemChart)
		//fmt.Printf("%s\n", p.Mem)
		system = append(system, p)
	}

	return system, nil
}

func GetWorker(coin string) ([]WorkerChart, error) {
	txn := Chart.Txn(false)
	defer txn.Abort()

	if coin == "ETH" {
		var worker []WorkerChart
		it, err := txn.Get("chart_ETH", "id")
		if err != nil {
			return nil, err
		}

		for obj := it.Next(); obj != nil; obj = it.Next() {
			p := obj.(WorkerChart)

			worker = append(worker, p)
		}

		return worker, nil
	} else {
		var worker []WorkerChart
		it, err := txn.Get("chart_ETC", "id")
		if err != nil {
			return nil, err
		}

		for obj := it.Next(); obj != nil; obj = it.Next() {
			p := obj.(WorkerChart)

			worker = append(worker, p)
		}

		return worker, nil
	}

}

func InsertTest() {
	txn := Chart.Txn(true)

	// Insert some people
	people := []*Miner{
		&Miner{global.Worker{
			Id:            "1231232132",
			Worker_name:   "t1",
			Wallet:        "0x123",
			Worker_idx:    100,
			Worker_share:  10,
			Worker_reject: 2,
		}},
		&Miner{global.Worker{
			Id:            "1231232132",
			Worker_name:   "t1",
			Wallet:        "0x123",
			Worker_idx:    100,
			Worker_share:  10,
			Worker_reject: 2,
		}},
		&Miner{global.Worker{
			Id:            "1231232133",
			Worker_name:   "t1",
			Wallet:        "0x123",
			Worker_idx:    100,
			Worker_share:  10,
			Worker_reject: 2,
		}},
		&Miner{global.Worker{
			Id:            "1231232134",
			Worker_name:   "t1",
			Wallet:        "0x123",
			Worker_idx:    100,
			Worker_share:  10,
			Worker_reject: 2,
		}},
	}
	for _, p := range people {
		if err := txn.Insert("miners", p); err != nil {
			panic(err)
		}
	}

	// Commit the transaction
	txn.Commit()
}

func ReadMiners() {

	txn := Chart.Txn(false)
	defer txn.Abort()
	// Lookup by email
	raw, err := txn.First("miners", "id", "1231232133")
	if err != nil {
		panic(err)
	}

	// Say hi!
	fmt.Printf("Hello %s!\n", raw.(*Miner).Worker_name)

	// List all the people
	it, err := txn.Get("miners", "id")
	if err != nil {
		panic(err)
	}

	for obj := it.Next(); obj != nil; obj = it.Next() {
		p := obj.(*Miner)
		fmt.Printf("  %s\n", p.Worker_name)
	}
}
