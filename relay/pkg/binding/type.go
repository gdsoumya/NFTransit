package binding

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type IEVM interface {
	Address() common.Address
	GetCurrentBlock() (*big.Int, error)
	UserTokens(contractAddr string, user common.Address) (map[uint64]string, error)
	AccountBalance(accountAddr string, self bool) (*big.Int, error)
	MintNonce(contractAddr string) (*big.Int, error)
	BurnNonce(contractAddr string) (*big.Int, error)
	GetTx(_txHash string) (*TxResponse, error)
	BurnLog(contractAddr string, _txHash string) (*BurnLogDetails, error)
	Mint(contractAddr string, _uris []string, _tos []common.Address, signature []byte) (*TxResponse, error)
	Burn(contractAddr string, ids []*big.Int) (*TxResponse, error)
	DeployToken(_admin common.Address, _name string, _version string) (string, error)
	GenerateMintEIP712Hash(contractAddr string, _uris []string, _tos []common.Address, nonce *big.Int) ([32]byte, error)
	GenerateBurnEIP712Hash(contractAddr string, nonce *big.Int) ([32]byte, error)
	VerifyEIP712Signature(hash common.Hash, _signature, _address string) (bool, error)
	Sign(hash common.Hash) ([]byte, error)
}

type BurnLogDetails struct {
	Ids         []*big.Int
	BurnNonce   *big.Int
	User        common.Address
	BlockNumber *big.Int
	Success     bool
}

type TxResponse struct {
	Cost        *big.Int
	TxHash      string
	Success     bool
	BlockNumber *big.Int
}
