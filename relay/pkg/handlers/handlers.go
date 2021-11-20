package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"strings"
	"time"

	"github.com/gdsoumya/nftransit/relay/pkg/binding"
	"github.com/gdsoumya/nftransit/relay/pkg/database"
	"github.com/gdsoumya/nftransit/relay/pkg/queue"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Handler struct {
	DB        database.IDatabase
	QClient   queue.IQueue
	Logger    *zap.Logger
	EvmClient binding.IEVM
}

func (h *Handler) QueueMint(c *gin.Context) {
	var data MintPayload
	err := c.BindJSON(&data)
	if err != nil {
		h.Logger.Debug("failed to parse request", zap.Error(err))
		c.JSON(400, gin.H{"error": fmt.Sprintf("failed to parse request, error: %v", err.Error())})
		return
	}
	uris, err := json.Marshal(data.URIs)
	if err != nil {
		h.Logger.Debug("failed to marshal uris", zap.Error(err), zap.Any("request", data))
		c.JSON(400, gin.H{"error": fmt.Sprintf("failed to parse request, error: %v", err.Error())})
		return
	}
	toAddrs, err := json.Marshal(data.ToAddrs)
	if err != nil {
		h.Logger.Debug("failed to marshal toAddrs", zap.Error(err), zap.Any("request", data))
		c.JSON(400, gin.H{"error": fmt.Sprintf("failed to parse request, error: %v", err.Error())})
		return
	}
	query, err := json.Marshal(MintQuery{Contract: data.Contract, Counter: data.Counter})
	if err != nil {
		h.Logger.Debug("failed to marshal mint query", zap.Error(err), zap.Any("request", data))
		c.JSON(400, gin.H{"error": fmt.Sprintf("failed to parse request, error: %v", err.Error())})
		return
	}

	nonce, err := h.EvmClient.MintNonce(data.Contract)
	if err != nil {
		h.Logger.Debug("failed to query mint nonce", zap.Error(err), zap.Any("request", data))
		c.JSON(400, gin.H{"error": fmt.Sprintf("failed to parse request, error: %v", err.Error())})
		return
	}
	if nonce.Cmp(new(big.Int).SetUint64(*data.Counter)) != 0 {
		h.Logger.Debug("mint nonce mismatch", zap.Any("expected", nonce), zap.Any("request", data))
		c.JSON(400, gin.H{"error": fmt.Sprintf("mint nonce mismatch")})
		return
	}

	// insert into db, unique key (contract + counter), and push to queue
	status := 500
	err = h.DB.Transaction(func(tx *gorm.DB) error {
		res, err := h.DB.Find(nil, &database.MintRequest{Contract: data.Contract, Counter: *data.Counter})
		if err != nil {
			return err
		}
		entry := database.MintRequest{
			CreatedAt: time.Now(),
			Contract:  data.Contract,
			Counter:   *data.Counter,
			Sender:    data.Sender,
			Status:    database.Pending,
			URIs:      string(uris),
			ToAddrs:   string(toAddrs),
			Signature: data.Signature,
		}
		if len(res) != 0 {
			if res[0].Status == database.Failed {
				if err := h.DB.Update(tx, &entry, true); err != nil {
					return err
				}
			} else {
				status = 400
				return fmt.Errorf("mint request with same counter already in queue")
			}
		} else if err := h.DB.Insert(tx, &entry); err != nil {
			return err
		}
		// delay message by 10 secs for the db tx to commit
		if err = h.QClient.Publish(string(query), 5000); err != nil {
			return err
		}
		if err = h.QClient.ConfirmPublish(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		h.Logger.Debug("failed mint request queuing", zap.Error(err), zap.Any("request", data))
		if status == 500 {
			c.JSON(status, gin.H{"error": fmt.Sprintf("failed to queue mint request")})
		} else {
			c.JSON(status, gin.H{"error": fmt.Sprintf("failed to queue mint request, error=%v", err)})
		}
		return
	}
	c.JSON(200, gin.H{"message": "mint request queued"})
}

func (h *Handler) QueryMintStatus(c *gin.Context) {
	// query db to see status, iff tx completed tx hash should be available
	var data MintQuery
	err := c.BindJSON(&data)
	if err != nil {
		h.Logger.Debug("failed to parse request", zap.Error(err))
		c.JSON(400, gin.H{"error": fmt.Sprintf("failed to parse request, error: %v", err.Error())})
		return
	}
	res, err := h.DB.Find(nil, &database.MintRequest{Contract: data.Contract, Counter: *data.Counter, Sender: data.Sender})
	if err != nil {
		h.Logger.Debug("failed get mint request", zap.Error(err), zap.Any("request", data))
		c.JSON(400, gin.H{"error": fmt.Sprintf("failed to fetch mint request")})
		return
	}
	if len(res) != 1 {
		h.Logger.Debug("expected len mismatch", zap.Any("expected", 1), zap.Any("got", len(res)), zap.Any("request", data))
		c.JSON(400, gin.H{"error": fmt.Sprintf("failed to fetch mint request")})
		return
	}

	// query current block number for confirmations
	blockNum := res[0].BlockNumber
	if res[0].TxHash != nil {
		resp, err := h.EvmClient.GetCurrentBlock()
		if err != nil {
			h.Logger.Debug("failed query tx", zap.Error(err), zap.Any("tx", resp))
			c.JSON(400, gin.H{"error": fmt.Sprintf("failed to fetch mint request")})
			return
		}
		blockNum = resp.Uint64()
	}

	var uris, toAddrs []string
	err = json.Unmarshal([]byte(res[0].URIs), &uris)
	if err != nil {
		h.Logger.Debug("failed parse uris", zap.Error(err), zap.Any("request", data))
		c.JSON(500, gin.H{"error": fmt.Sprintf("failed to fetch mint request")})
		return
	}
	err = json.Unmarshal([]byte(res[0].ToAddrs), &toAddrs)
	if err != nil {
		h.Logger.Debug("failed get parse toAddrs", zap.Error(err), zap.Any("request", data))
		c.JSON(500, gin.H{"error": fmt.Sprintf("failed to fetch mint request")})
		return
	}
	c.JSON(200, MintResponse{
		CreatedAt:     res[0].CreatedAt,
		Contract:      res[0].Contract,
		Counter:       res[0].Counter,
		Confirmations: blockNum - res[0].BlockNumber,
		BlockNumber:   res[0].BlockNumber,
		URIs:          uris,
		ToAddrs:       toAddrs,
		Signature:     res[0].Signature,
		TxHash:        res[0].TxHash,
		Status:        res[0].Status,
		Error:         res[0].Error,
	})
}

func (h *Handler) GetBurnTx(c *gin.Context) {
	// query node to get burn tx log
	var data BurnTxQuery
	err := c.BindJSON(&data)
	if err != nil {
		h.Logger.Debug("failed to parse request", zap.Error(err))
		c.JSON(400, gin.H{"error": fmt.Sprintf("failed to parse request, error: %v", err.Error())})
		return
	}
	details, err := h.EvmClient.BurnLog(data.Contract, data.TxHash)
	if err != nil {
		h.Logger.Debug("failed get get burn tx logs", zap.Error(err), zap.Any("request", data))
		c.JSON(400, gin.H{"error": fmt.Sprintf("failed get get burn tx logs")})
		return
	}
	ids := []uint64{}
	for _, id := range details.Ids {
		ids = append(ids, id.Uint64())
	}
	c.JSON(200, BurnTxResponse{
		Ids:         ids,
		BurnNonce:   details.BurnNonce.Uint64(),
		User:        details.User.String(),
		BlockNumber: details.BlockNumber.Uint64(),
		Success:     details.Success,
	})
}

func (h *Handler) VerifyBurnTx(c *gin.Context) {
	var data VerifyBurnQuery
	err := c.BindJSON(&data)
	if err != nil {
		h.Logger.Debug("failed to parse request", zap.Error(err))
		c.JSON(400, gin.H{"error": fmt.Sprintf("failed to parse request, error: %v", err.Error())})
		return
	}
	details, err := h.EvmClient.BurnLog(data.Contract, data.TxHash)
	if err != nil {
		h.Logger.Debug("failed get get burn tx logs", zap.Error(err), zap.Any("request", data))
		c.JSON(400, gin.H{"error": fmt.Sprintf("failed get get burn tx logs")})
		return
	}
	hash, err := h.EvmClient.GenerateBurnEIP712Hash(data.Contract, details.BurnNonce)
	if err != nil {
		h.Logger.Debug("failed to generate burn hash", zap.Error(err), zap.Any("request", data))
		c.JSON(400, gin.H{"error": fmt.Sprintf("failed to generate burn hash")})
		return
	}
	valid, err := h.EvmClient.VerifyEIP712Signature(hash, strings.Replace(data.Signature, "0x", "", 1), details.User.String())
	if err != nil {
		h.Logger.Debug("failed to verify burn signature", zap.Error(err), zap.Any("request", data))
		c.JSON(400, gin.H{"error": fmt.Sprintf("failed to verify burn signature")})
		return
	}
	if !valid {
		h.Logger.Debug("invalid burn signature", zap.Any("request", data),zap.Any("burn_details",details))
		c.JSON(400, gin.H{"error": fmt.Sprintf("invalid burn signature")})
		return
	}
	ids := []uint64{}
	for _, id := range details.Ids {
		ids = append(ids, id.Uint64())
	}
	c.JSON(200, VerifyBurnResponse{
		BurnTxResponse{
			Ids:         ids,
			BurnNonce:   details.BurnNonce.Uint64(),
			User:        details.User.String(),
			BlockNumber: details.BlockNumber.Uint64(),
			Success:     details.Success},
		true,
	})
}

func (h *Handler) GetUserTokens(c *gin.Context) {
	var data UserTokenQuery
	err := c.BindJSON(&data)
	if err != nil {
		h.Logger.Debug("failed to parse request", zap.Error(err))
		c.JSON(400, gin.H{"error": fmt.Sprintf("failed to parse request, error: %v", err.Error())})
		return
	}
	details, err := h.EvmClient.UserTokens(data.Contract, common.HexToAddress(data.Address))
	if err != nil {
		h.Logger.Debug("failed get get user tokens", zap.Error(err), zap.Any("request", data))
		c.JSON(400, gin.H{"error": fmt.Sprintf("failed get get user tokens")})
		return
	}
	c.JSON(200, details)
}
