//go:generate abigen --abi ../../../contracts/token.abi --bin ../../../contracts/token.bin  --pkg binding --out binding.go
package binding

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"unsafe"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

type EVMClient struct {
	logger  *zap.Logger
	conn    *ethclient.Client
	privKey *ecdsa.PrivateKey
	chainID *big.Int
}

func NewEVMClient(logger *zap.Logger, rpcURL string, privKey string, chainID *big.Int) (*EVMClient, error) {
	conn, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to node, error:%w", err)
	}
	pk, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return nil, fmt.Errorf("invalid priv key, error:%w", err)
	}
	return &EVMClient{
		logger:  logger,
		conn:    conn,
		privKey: pk,
		chainID: chainID,
	}, nil
}

func (c *EVMClient) Address() common.Address {
	return crypto.PubkeyToAddress(c.privKey.PublicKey)
}

func (c *EVMClient) AccountBalance(accountAddr string, self bool) (*big.Int, error) {
	if self {
		return c.conn.BalanceAt(context.Background(), crypto.PubkeyToAddress(c.privKey.PublicKey), nil)
	}
	return c.conn.BalanceAt(context.Background(), common.HexToAddress(accountAddr), nil)
}

func (c *EVMClient) MintNonce(contractAddr string) (*big.Int, error) {
	contract, err := NewBinding(common.HexToAddress(contractAddr), c.conn)
	if err != nil {
		return nil, fmt.Errorf("failed to bind to contract, error:%w", err)
	}
	return contract.MintNonce(nil)
}

func (c *EVMClient) BurnNonce(contractAddr string) (*big.Int, error) {
	contract, err := NewBinding(common.HexToAddress(contractAddr), c.conn)
	if err != nil {
		return nil, fmt.Errorf("failed to bind to contract, error:%w", err)
	}
	return contract.BurnNonce(nil)
}

func (c *EVMClient) UserTokens(contractAddr string, user common.Address) (map[uint64]string, error) {
	contract, err := NewBinding(common.HexToAddress(contractAddr), c.conn)
	if err != nil {
		return nil, fmt.Errorf("failed to bind to contract, error:%w", err)
	}
	tokens, err := contract.UserTokens(nil, user)
	if err != nil {
		return nil, fmt.Errorf("failed to to get user tokens, error:%w", err)
	}
	userTokens := map[uint64]string{}
	for i, token := range tokens {
		if token != "" {
			userTokens[uint64(i)] = token
		}
	}
	return userTokens, nil
}

func (c *EVMClient) BurnLog(contractAddr string, _txHash string) (*BurnLogDetails, error) {
	contract, err := NewBinding(common.HexToAddress(contractAddr), c.conn)
	if err != nil {
		return nil, fmt.Errorf("failed to bind to contract, error:%w", err)
	}
	txHash := common.HexToHash(_txHash)
	receipt, err := c.conn.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get tx reciept, error:%w", err)
	}
	var burnLog *BindingTokenBurn = nil
	for i := range receipt.Logs {
		burnLog, err = contract.ParseTokenBurn(*receipt.Logs[i])
		if err == nil {
			break
		}
	}
	if burnLog == nil {
		return nil, fmt.Errorf("failed to get burn log from tx")
	}

	return &BurnLogDetails{
		Ids:         burnLog.Ids,
		BurnNonce:   burnLog.BurnNonce,
		User:        burnLog.User,
		BlockNumber: receipt.BlockNumber,
		Success:     receipt.Status == 1,
	}, nil
}

func (c *EVMClient) GetTx(_txHash string) (*TxResponse, error) {
	txHash := common.HexToHash(_txHash)
	receipt, err := c.conn.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get tx reciept, error:%w", err)
	}
	return &TxResponse{
		BlockNumber: receipt.BlockNumber,
		Success:     receipt.Status == 1,
	}, nil
}

func (c *EVMClient) GetCurrentBlock() (*big.Int, error) {
	head, err := c.conn.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest block, error:%w", err)
	}
	return head.Number, nil
}

func (c *EVMClient) Mint(contractAddr string, _uris []string, _tos []common.Address, signature []byte) (*TxResponse, error) {
	contract, err := NewBinding(common.HexToAddress(contractAddr), c.conn)
	if err != nil {
		return nil, fmt.Errorf("failed to bind to contract, error:%w", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(c.privKey, c.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create authorized transactor, error:%w", err)
	}
	log.Print("from",auth.From)
	tx, err := contract.Mint(auth, _uris, _tos, signature)
	if err != nil {
		return nil, fmt.Errorf("failed to create mint tx, error:%w", err)
	}
	c.logger.Debug("mint tx sent, waiting for confirmation", zap.Any("txHash", tx.Hash().String()), zap.Any("uris", _uris), zap.Any("tos", _tos), zap.Any("signature", signature))
	receipt, err := bind.WaitMined(context.Background(), c.conn, tx)
	if err != nil {
		return nil, fmt.Errorf("tx not included in block, error:%w", err)
	}
	cost := new(big.Int).SetUint64(receipt.GasUsed)
	cost.Mul(cost, tx.GasPrice())

	return &TxResponse{
		Cost:        cost,
		TxHash:      tx.Hash().String(),
		Success:     receipt.Status == 1,
		BlockNumber: receipt.BlockNumber,
	}, nil
}

func (c *EVMClient) Burn(contractAddr string, ids []*big.Int) (*TxResponse, error) {
	contract, err := NewBinding(common.HexToAddress(contractAddr), c.conn)
	if err != nil {
		return nil, fmt.Errorf("failed to bind to contract, error:%w", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(c.privKey, c.chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create authorized transactor, error:%w", err)
	}

	tx, err := contract.Burn(auth, ids)
	if err != nil {
		return nil, fmt.Errorf("failed to create burn tx, error:%w", err)
	}
	c.logger.Debug("burn tx sent, waiting for confirmation", zap.Any("ids", ids), zap.Any("txHash", tx.Hash().String()))
	receipt, err := bind.WaitMined(context.Background(), c.conn, tx)
	if err != nil {
		return nil, fmt.Errorf("tx not included in block, error:%w", err)
	}
	cost := new(big.Int).SetUint64(receipt.GasUsed)
	cost.Mul(cost, tx.GasPrice())

	return &TxResponse{
		Cost:        cost,
		TxHash:      tx.Hash().String(),
		Success:     receipt.Status == 1,
		BlockNumber: receipt.BlockNumber,
	}, nil
}

func (c *EVMClient) DeployToken(_admin common.Address, _name string, _version string) (string, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(c.privKey, c.chainID)
	if err != nil {
		return "", fmt.Errorf("failed to create authorized transactor, error:%w", err)
	}

	contractAddr, tx, _, err := DeployBinding(auth, c.conn, _admin, _name, _version)
	if err != nil {
		return "", fmt.Errorf("failed to deploy contract, error:%w", err)
	}
	c.logger.Debug("deploy tx sent, waiting for confirmation", zap.Any("contractAddr", contractAddr), zap.Any("txHash", tx.Hash().String()))
	addr, err := bind.WaitDeployed(context.Background(), c.conn, tx)
	if err != nil {
		return "", fmt.Errorf("tx not included in block, error:%w", err)
	}

	return addr.String(), nil
}

func (c *EVMClient) GenerateMintEIP712Hash(contractAddr string, _uris []string, _tos []common.Address, nonce *big.Int) ([32]byte, error) {
	contract, err := NewBinding(common.HexToAddress(contractAddr), c.conn)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to bind to contract, error:%w", err)
	}
	return contract.HashMintRequest(nil, _uris, _tos, nonce)
}

func (c *EVMClient) GenerateBurnEIP712Hash(contractAddr string, nonce *big.Int) ([32]byte, error) {
	contract, err := NewBinding(common.HexToAddress(contractAddr), c.conn)
	if err != nil {
		return [32]byte{}, fmt.Errorf("failed to bind to contract, error:%w", err)
	}
	return contract.HashBurnRequest(nil, nonce)
}

func (c *EVMClient) VerifyEIP712Signature(hash common.Hash, _signature, _address string) (bool, error) {
	signature, _ := hex.DecodeString(_signature)
	if len(signature) != 65 {
		return false, fmt.Errorf("invalid signature length: %d", len(signature))
	}

	if signature[64] != 27 && signature[64] != 28 {
		return false, fmt.Errorf("invalid recovery id: %d", signature[64])
	}
	signature[64] -= 27

	pubKeyRaw, err := crypto.Ecrecover(hash.Bytes(), signature)
	if err != nil {
		return false, fmt.Errorf("invalid signature: %s", err.Error())
	}

	pubKey, err := crypto.UnmarshalPubkey(pubKeyRaw)
	if err != nil {
		return false, err
	}

	address := common.HexToAddress(_address)
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	return bytes.Equal(address.Bytes(), recoveredAddr.Bytes()), nil
}

func (c *EVMClient) Sign(hash common.Hash) ([]byte, error) {
	sig, err := crypto.Sign(hash.Bytes(), c.privKey)
	if err != nil {
		return nil, err
	}
	if sig[64] < 27 {
		sig[64] += 27
	}
	return sig, nil
}

func byte32(b []byte) [32]byte {
	return *(*[32]byte)(unsafe.Pointer(&b))
}
