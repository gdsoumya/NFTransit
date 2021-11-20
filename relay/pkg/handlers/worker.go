package handlers

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gdsoumya/nftransit/relay/pkg/binding"
	"github.com/gdsoumya/nftransit/relay/pkg/util"

	"github.com/gdsoumya/nftransit/relay/pkg/database"
	"github.com/gdsoumya/nftransit/relay/pkg/env"
	"github.com/gdsoumya/nftransit/relay/pkg/queue"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	MaxRetries = 2
)

func (h *Handler) StartWorkers(envConfig env.EnvData) {
	for i, _ := range envConfig.PrivateKeys {
		go h.Work(i, envConfig)
	}
}

func (h *Handler) Work(i int, envConfig env.EnvData) {
	h.Logger.Info("started mint worker", zap.Int("id", i), zap.String("pk", envConfig.PrivateKeys[i][0:5]+"..."))
	client, err := queue.InitQueue(h.Logger, &envConfig, false)
	if err != nil {
		h.Logger.Fatal("failed to get consumer queue client", zap.Int("id", i), zap.Error(err))
	}
	retryClient, err := queue.InitQueue(h.Logger, &envConfig, true)
	if err != nil {
		h.Logger.Fatal("failed to get consumer queue client", zap.Int("id", i), zap.Error(err))
	}

	//create chain clients
	evmClient, err := binding.NewEVMClient(h.Logger, envConfig.NodeURL, envConfig.PrivateKeys[i], big.NewInt(envConfig.ChainID))
	if err != nil {
		h.Logger.Fatal("failed to get chain client", zap.Int("id", i), zap.Error(err))
	}

	defer client.Close()
	defer retryClient.Close()

	msgs, err := client.Consume()

	if err != nil {
		h.Logger.Fatal("failed get msgs for consumer", zap.Int("id", i), zap.Error(err))
	}

	for msg := range msgs {
		h.Logger.Debug("handling mint request", zap.Int("id", i), zap.String("msg", string(msg.Body)))
		data := MintQuery{}
		if err = json.Unmarshal(msg.Body, &data); err != nil {
			h.Logger.Error("failed to unmarshal msg", zap.Int("id", i), zap.Error(err), zap.String("msg", string(msg.Body)))
			continue
		}

		err = h.DB.Transaction(func(tx *gorm.DB) error {
			// find item in db
			res, err := h.DB.Find(tx, &database.MintRequest{Contract: data.Contract, Counter: *data.Counter})
			if err != nil {
				return err
			}
			if len(res) != 1 {
				return fmt.Errorf("failed to fetch mint request from db")
			}
			item := res[0]
			if item.Status != database.Pending {
				return nil
			}
			// perform tx
			uris := []string{}
			err = json.Unmarshal([]byte(item.URIs), &uris)
			if err != nil {
				return fmt.Errorf("failed to parse uris, error:%w", err)
			}
			_toAddrs := []string{}
			err = json.Unmarshal([]byte(item.ToAddrs), &_toAddrs)
			if err != nil {
				return fmt.Errorf("failed to parse to addrs, error:%w", err)
			}
			toAddrs := []common.Address{}
			for _, addr := range _toAddrs {
				toAddrs = append(toAddrs, common.HexToAddress(addr))
			}
			sig, err := hex.DecodeString(strings.Replace(item.Signature, "0x", "", 1))
			if err != nil {
				return fmt.Errorf("failed to parse sig, error:%w", err)
			}
			resp, err := evmClient.Mint(item.Contract, uris, toAddrs, sig)
			// if tx reverted set error
			if err != nil && strings.Contains(strings.ToLower(err.Error()), "vm exception") && strings.Contains(strings.ToLower(err.Error()), "revert") {
				item.Status = database.Failed
				item.Error = util.StringToPtr(err.Error())
			} else if err != nil {
				return fmt.Errorf("mint tx failed, error:%w", err)
			} else {
				// update status in db
				if resp.Success {
					item.Status = database.Completed
				} else {
					item.Status = database.Failed
				}
				item.TxHash = util.StringToPtr(resp.TxHash)
				item.BlockNumber = resp.BlockNumber.Uint64()
				item.Cost = resp.Cost.Uint64()
			}
			// update in db
			err = h.DB.Update(tx, &item, false)
			if err != nil {
				return err
			}
			return nil
		})

		// if above tx fails and max retries have been reached set to failed
		if err != nil && data.Retries == MaxRetries {
			err = h.DB.Transaction(func(tx *gorm.DB) error {
				// find item in db
				res, err1 := h.DB.Find(tx, &database.MintRequest{Contract: data.Contract, Counter: *data.Counter})
				if err1 != nil {
					return err1
				}
				if len(res) != 1 {
					return fmt.Errorf("failed to fetch mint request from db")
				}
				item := res[0]
				if item.Status != database.Pending {
					return nil
				}
				// update failed status in db
				errMsg := fmt.Sprintf("Failed to perform mint request after %v tries, error:%v", MaxRetries, err)
				item.Status = database.Failed
				item.Error = &errMsg
				err = h.DB.Update(tx, &item, false)
				if err != nil {
					return err
				}
				return nil
			})
		}

		// if either of the above fails keep retrying until no more tries left
		if err != nil {
			if data.Retries >= MaxRetries {
				h.Logger.Error("failed handle message, max retries over dropping request", zap.Int("id", i), zap.Error(err), zap.String("msg", string(msg.Body)))
			} else {
				h.Logger.Error("failed handle message", zap.Int("id", i), zap.Error(err), zap.String("msg", string(msg.Body)))
				data.Retries += 1
				retryMsg, err := json.Marshal(data)
				if err != nil {
					h.Logger.Error("failed marshal retry msg", zap.Int("id", i), zap.Error(err), zap.Any("msg", retryMsg))
				}
				// retry 15 secs later
				if err = retryClient.Publish(string(retryMsg), 5000); err != nil {
					h.Logger.Error("failed publish retry message", zap.Int("id", i), zap.Error(err), zap.String("msg", string(retryMsg)))
				}
				h.Logger.Debug("requeue")
				if err = retryClient.ConfirmPublish(); err != nil {
					h.Logger.Error("failed confirm retry message publish", zap.Int("id", i), zap.Error(err), zap.String("msg", string(retryMsg)))
				}
				h.Logger.Debug("confirm requeue")
			}
		}
		if err = client.AckDelivery(&msg, false); err != nil {
			h.Logger.Error("failed ack message", zap.Int("id", i), zap.Error(err), zap.String("msg", string(msg.Body)))
		} else {
			h.Logger.Debug("original removed acked")
		}
	}
}
