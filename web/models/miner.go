package models

import (
	"fmt"
	"miner_proxy/pack"

	memdb "github.com/hashicorp/go-memdb"
)

var MinerDB *memdb.MemDB

func init() {
	schema := &memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"miners": &memdb.TableSchema{
				Name: "miners",
				Indexes: map[string]*memdb.IndexSchema{
					"id": &memdb.IndexSchema{
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "Id"},
					},
					"wallet": &memdb.IndexSchema{
						Name:    "wallet",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Wallet"},
					},
					"worker": &memdb.IndexSchema{
						Name:    "worker",
						Unique:  false,
						Indexer: &memdb.StringFieldIndex{Field: "Worker_name"},
					},
				},
			},
		},
	}

	// Create a new data base
	var err error
	MinerDB, err = memdb.NewMemDB(schema)
	if err != nil {
		panic(err)
	}
}

type Miner struct {
	pack.Worker
}

func NewMiner(w pack.Worker) (*Miner, error) {
	return &Miner{w}, nil
}

func InsertTest() {
	txn := MinerDB.Txn(true)

	// Insert some people
	people := []*Miner{
		&Miner{pack.Worker{
			Id:            "1231232132",
			Worker_name:   "t1",
			Wallet:        "0x123",
			Worker_idx:    100,
			Worker_share:  10,
			Worker_reject: 2,
		}},
		&Miner{pack.Worker{
			Id:            "1231232132",
			Worker_name:   "t1",
			Wallet:        "0x123",
			Worker_idx:    100,
			Worker_share:  10,
			Worker_reject: 2,
		}},
		&Miner{pack.Worker{
			Id:            "1231232133",
			Worker_name:   "t1",
			Wallet:        "0x123",
			Worker_idx:    100,
			Worker_share:  10,
			Worker_reject: 2,
		}},
		&Miner{pack.Worker{
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

	txn := MinerDB.Txn(false)
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

	fmt.Println("All the people:")
	for obj := it.Next(); obj != nil; obj = it.Next() {
		p := obj.(*Miner)
		fmt.Printf("  %s\n", p.Worker_name)
	}

	// Range scan over people with ages between 25 and 35 inclusive
	// it, err = txn.LowerBound("person", "age", 25)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("People aged 25 - 35:")
	// for obj := it.Next(); obj != nil; obj = it.Next() {
	// 	p := obj.(*Person)
	// 	if p.Age > 35 {
	// 		break
	// 	}
	// 	fmt.Printf("  %s is aged %d\n", p.Name, p.Age)
	// }
}