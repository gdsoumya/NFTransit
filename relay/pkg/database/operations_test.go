package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/gdsoumya/nftransit/relay/pkg/util"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type MetadataTest struct {
	Uri    string `json:"uri"`
	ToAddr string `json:"to"`
}

func TestDB(t *testing.T) {

	// setup db connection
	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("failed to create logger, error: %v", err)
	}
	db, err := NewDatabaseClient(logger, &DBConfig{
		User:     "postgres",
		Password: "Pass2020",
		DBName:   "nftransit_test",
		Host:     "127.0.0.1",
		Port:     5432,
	}, true)
	if err != nil {
		t.Fatalf("failed to connect to db, error: %v", err)
	}

	t.Run("clean start", func(t1 *testing.T) {
		if err := db.DropTable(); err != nil {
			t1.Fatalf("failed to drop mint_requests table, err:%v", err)
		}
	})

	t.Run("run dm migration", func(t1 *testing.T) {
		if err := db.db.AutoMigrate(&MintRequest{}); err != nil {
			t1.Fatalf("failed to migrate mint_requests table, err:%v", err)
		}
	})

	t.Run("insert", func(t1 *testing.T) {
		uri, err := json.Marshal([]string{"uri1", "uri2"})
		if err != nil {
			t1.Fatalf("faile to marshal metadata, err:%v", err)
		}
		addrs, err := json.Marshal([]string{"0x123", "0x1234"})
		if err != nil {
			t1.Fatalf("faile to marshal metadata, err:%v", err)
		}
		if err = db.Insert(nil, &MintRequest{
			Contract:  "0x12435445",
			Counter:   1,
			URIs:      string(uri),
			ToAddrs:   string(addrs),
			Signature: "sig1",
			Sender:    "0x243534",
			Status:    Pending,
		}); err != nil {
			t1.Fatalf("failed to insert in DB, err:%v", err)
		}

		if err = db.Insert(nil, &MintRequest{
			Contract:  "0x12435445",
			Counter:   1,
			URIs:      string(uri),
			ToAddrs:   string(addrs),
			Signature: "sig1",
			Sender:    "0x243534",
			Status:    Pending,
		}); err == nil {
			t1.Fatalf("expected to fail but passed")
		}

	})

	t.Run("query", func(t1 *testing.T) {
		expected := MintRequest{
			Contract:  "0x12435445",
			Counter:   1,
			URIs:      `["uri1","uri2"]`,
			ToAddrs:   `["0x123","0x1234"]`,
			Signature: "sig1",
			Status:    Pending,
			Sender:    "0x243534",
		}
		cases := []MintRequest{{
			Contract: "0x12435445",
			Counter:  1},
			{
				Contract: "0x12435445",
			}}
		for _, c := range cases {
			data, err := db.Find(nil, &c)
			if err != nil {
				t1.Fatalf("failed to query in DB, err:%v", err)
			}
			if len(data) != 1 {
				t1.Fatalf("recieved elements of len=%v expected %v", len(data), 1)
			}
			if data[0].Contract != expected.Contract {
				t1.Fatalf("recieved elements don't match expected value expected=%v, recieved=%v", expected, data[0])
			}
			if data[0].Counter != expected.Counter {
				t1.Fatalf("recieved elements don't match expected value expected=%v, recieved=%v", expected, data[0])
			}
			if data[0].URIs != expected.URIs {
				t1.Fatalf("recieved elements don't match expected value expected=%v, recieved=%v", expected, data[0])
			}
			if data[0].ToAddrs != expected.ToAddrs {
				t1.Fatalf("recieved elements don't match expected value expected=%v, recieved=%v", expected, data[0])
			}
			if data[0].Status != expected.Status {
				t1.Fatalf("recieved elements don't match expected value expected=%v, recieved=%v", expected, data[0])
			}
			if data[0].Sender != expected.Sender {
				t1.Fatalf("recieved elements don't match expected value expected=%v, recieved=%v", expected, data[0])
			}
			if data[0].Signature != expected.Signature {
				t1.Fatalf("recieved elements don't match expected value expected=%v, recieved=%v", expected, data[0])
			}
		}
	})

	t.Run("update", func(t1 *testing.T) {
		expected := MintRequest{
			Contract:  "0x12435445",
			Counter:   1,
			URIs:      `["uri1","uri2"]`,
			ToAddrs:   `["0x123","0x1234"]`,
			Signature: "sig1",
			Status:    Completed,
			Sender:    "0x243534",
			TxHash:    util.StringToPtr("0x5646546757"),
		}
		err = db.Update(nil, &MintRequest{
			Contract: "0x12435445",
			Counter:  1,
			Status:   Completed,
			TxHash:   util.StringToPtr("0x5646546757"),
		}, false)
		if err != nil {
			t1.Fatalf("failed to update in DB, err:%v", err)
		}
		data, err := db.Find(nil, &MintRequest{
			Contract: "0x12435445",
			Counter:  1})
		if err != nil {
			t1.Fatalf("failed to query in DB, err:%v", err)
		}
		if len(data) != 1 {
			t1.Fatalf("recieved elements of len=%v expected %v", len(data), 1)
		}
		if data[0].Contract != expected.Contract {
			t1.Fatalf("recieved elements don't match expected value expected=%v, recieved=%v", expected, data[0])
		}
		if data[0].Counter != expected.Counter {
			t1.Fatalf("recieved elements don't match expected value expected=%v, recieved=%v", expected, data[0])
		}
		if data[0].URIs != expected.URIs {
			t1.Fatalf("recieved elements don't match expected value expected=%v, recieved=%v", expected, data[0])
		}
		if data[0].ToAddrs != expected.ToAddrs {
			t1.Fatalf("recieved elements don't match expected value expected=%v, recieved=%v", expected, data[0])
		}
		if data[0].Status != expected.Status {
			t1.Fatalf("recieved elements don't match expected value expected=%v, recieved=%v", expected, data[0])
		}
		if data[0].Sender != expected.Sender {
			t1.Fatalf("recieved elements don't match expected value expected=%v, recieved=%v", expected, data[0])
		}
		if *data[0].TxHash != *expected.TxHash {
			t1.Fatalf("recieved elements don't match expected value expected=%v, recieved=%v", expected, data[0])
		}
		if data[0].Signature != expected.Signature {
			t1.Fatalf("recieved elements don't match expected value expected=%v, recieved=%v", expected, data[0])
		}
	})

	t.Run("delete", func(t1 *testing.T) {
		err = db.Delete(nil, &MintRequest{
			Contract: "0x12435445",
			Counter:  1,
		})
		if err != nil {
			t1.Fatalf("failed to delete in DB, err:%v", err)
		}
		data, err := db.Find(nil, &MintRequest{
			Contract: "0x12435445",
			Counter:  1})
		if err != nil {
			t1.Fatalf("failed to query in DB, err:%v", err)
		}
		if len(data) != 0 {
			t1.Fatalf("delete failed, recieved elements of len=%v expected %v", len(data), 0)
		}
	})

	t.Run("transaction", func(t1 *testing.T) {
		expected := MintRequest{
			Contract:  "0x12435445",
			Counter:   1,
			URIs:      `["uri1","uri2"]`,
			ToAddrs:   `["0x123","0x1234"]`,
			Signature: "sig1",
			Status:    Completed,
			Sender:    "0x243534",
			TxHash:    util.StringToPtr("0x5646546757"),
		}
		err := db.Transaction(func(tx *gorm.DB) error {
			if err = db.Insert(tx, &expected); err != nil {
				return fmt.Errorf("failed to insert in DB, err:%v", err)
			}
			if err = db.Update(tx, &MintRequest{
				Contract: "0x12435445",
				Counter:  1,
				Status:   Pending,
				TxHash:   util.StringToPtr("0x5646546757"),
			}, false); err != nil {
				return fmt.Errorf("failed to update in DB, err:%v", err)
			}
			return errors.New("failed")
		})
		if err == nil {
			t1.Fatalf("tx should have failed but passed")
		}
		data, err := db.Find(nil, &MintRequest{
			Contract: "0x12435445",
			Counter:  1})
		if err != nil {
			t1.Fatalf("failed to query in DB, err:%v", err)
		}
		if len(data) != 0 {
			t1.Fatalf("tx passed, recieved elements of len=%v expected %v", len(data), 0)
		}

		err = db.Transaction(func(tx *gorm.DB) error {
			if err = db.Insert(tx, &expected); err != nil {
				return fmt.Errorf("failed to insert in DB, err:%v", err)
			}
			if err = db.Update(tx, &MintRequest{
				Contract: "0x12435445",
				Counter:  1,
				Status:   Pending,
				TxHash:   util.StringToPtr("0x5646546757"),
			}, false); err != nil {
				return fmt.Errorf("failed to update in DB, err:%v", err)
			}
			return nil
		})
		if err != nil {
			t1.Fatalf("tx should have passed but failed, err:%v", err)
		}
		data, err = db.Find(nil, &MintRequest{
			Contract: "0x12435445",
			Counter:  1})
		if err != nil {
			t1.Fatalf("failed to query in DB, err:%v", err)
		}
		if len(data) != 1 {
			t1.Fatalf("tx failed, recieved elements of len=%v expected %v", len(data), 0)
		}
		if data[0].Contract != expected.Contract {
			t1.Fatalf("recieved elements don't match expected value expected=%v, recieved=%v", expected, data[0])
		}
		if data[0].Counter != expected.Counter {
			t1.Fatalf("recieved elements don't match expected value expected=%v, recieved=%v", expected, data[0])
		}
		if data[0].URIs != expected.URIs {
			t1.Fatalf("recieved elements don't match expected value expected=%v, recieved=%v", expected, data[0])
		}
		if data[0].ToAddrs != expected.ToAddrs {
			t1.Fatalf("recieved elements don't match expected value expected=%v, recieved=%v", expected, data[0])
		}
		if data[0].Status != Pending {
			t1.Fatalf("recieved elements don't match expected value expected=%v, recieved=%v", expected, data[0])
		}
		if data[0].Sender != expected.Sender {
			t1.Fatalf("recieved elements don't match expected value expected=%v, recieved=%v", expected, data[0])
		}
		if *data[0].TxHash != *expected.TxHash {
			t1.Fatalf("recieved elements don't match expected value expected=%v, recieved=%v", expected, data[0])
		}
		if data[0].Signature != expected.Signature {
			t1.Fatalf("recieved elements don't match expected value expected=%v, recieved=%v", expected, data[0])
		}
	})
}
