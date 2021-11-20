package handlers

import (
	"time"

	"github.com/gdsoumya/nftransit/relay/pkg/database"
)

type Metadata struct {
	Uri    string `json:"uri" binding:"required"`
	ToAddr string `json:"to" binding:"required"`
}

type MintPayload struct {
	Contract  string   `json:"contract" binding:"required"`
	Counter   *uint64  `json:"counter" binding:"required"`
	URIs      []string `json:"uris" binding:"required"`
	ToAddrs   []string `json:"to_addrs" binding:"required"`
	Signature string   `json:"signature" binding:"required"`
	Sender    string   `json:"sender" binding:"required"`
}

type MintResponse struct {
	CreatedAt     time.Time           `json:"create_at"`
	Contract      string              `json:"contract"`
	Counter       uint64              `json:"counter"`
	Confirmations uint64              `json:"confirmations"`
	BlockNumber   uint64              `json:"block_number"`
	URIs          []string            `json:"uris"`
	ToAddrs       []string            `json:"to_addrs"`
	TxHash        *string             `json:"tx_hash"`
	Signature     string              `json:"signature"`
	Status        database.MintStatus `json:"status"`
	Error         *string             `json:"error"`
}

type MintQuery struct {
	Contract string  `json:"contract" binding:"required"`
	Counter  *uint64 `json:"counter" binding:"required"`
	Sender   string  `json:"sender" binding:"required"`
	Retries  uint64  `json:"retries"`
}

type BurnTxQuery struct {
	Contract string `json:"contract" binding:"required"`
	TxHash   string `json:"tx_hash" binding:"required"`
}

type VerifyBurnQuery struct {
	Contract  string `json:"contract" binding:"required"`
	TxHash    string `json:"tx_hash" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

type BurnTxResponse struct {
	Ids         []uint64 `json:"ids"`
	BurnNonce   uint64   `json:"nonce"`
	User        string   `json:"user"`
	BlockNumber uint64   `json:"block_number"`
	Success     bool     `json:"success"`
}

type VerifyBurnResponse struct {
	BurnTxResponse
	Valid bool `json:"valid"`
}

type UserTokenQuery struct {
	Contract string `json:"contract" binding:"required"`
	Address  string `json:"address" binding:"required"`
}
