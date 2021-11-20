package binding

import (
	"encoding/hex"
	"math/big"
	"os"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

func TestEVMClient(t *testing.T) {
	pkString := os.Getenv("PRIVATE_KEY1")
	signerString := os.Getenv("PRIVATE_KEY2")

	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("failed to create logger, error: %v", err)
	}

	// setup chain client
	client, err := NewEVMClient(logger, "http://127.0.0.1:7545", pkString, big.NewInt(75))
	if err != nil {
		t.Fatalf("failed to create chain client, error:%v", err)
	}

	// setup client for signer
	signerClient, err := NewEVMClient(logger, "http://127.0.0.1:7545", signerString, big.NewInt(75))
	if err != nil {
		t.Fatalf("failed to create chain client, error:%v", err)
	}

	var (
		contract  = ""
		name      = "NFTransit"
		version   = "0.1"
		mintSig   = []byte{}
		burnNonce = new(big.Int)
		burnTx    = ""
	)

	t.Run("deploy contract", func(t1 *testing.T) {
		contract, err = client.DeployToken(signerClient.Address(), name, version)
		if err != nil {
			t1.Fatalf("failed to deploy contract, error: %v", err)
		}
		t1.Logf("contract deployed at :%v", contract)
	})

	t.Run("mint directly", func(t1 *testing.T) {
		resp, err := signerClient.Mint(contract, []string{"uri0"}, []common.Address{client.Address()}, []byte{})
		if err != nil {
			t1.Fatalf("failed to make direct mint, error: %v", err)
		}
		if !resp.Success {
			t1.Fatalf("failed to make direct mint, error: tx failed")
		}
		t1.Logf("mint tx | cost=%v | block number=%v", resp.Cost, resp.BlockNumber)
	})

	t.Run("generate and verify mint signature", func(t1 *testing.T) {
		resp, err := signerClient.GenerateMintEIP712Hash(contract, []string{"uri1"}, []common.Address{client.Address()}, big.NewInt(1))
		if err != nil {
			t1.Fatalf("failed to generate mint eip712 hash, error: %v", err)
		}
		mintSig, err = signerClient.Sign(resp)
		if err != nil {
			t1.Fatalf("failed to sign mint eip712 hash, error: %v", err)
		}
		if valid, err := signerClient.VerifyEIP712Signature(resp, hex.EncodeToString(mintSig), signerClient.Address().String()); err != nil || !valid {
			t1.Fatalf("failed to generate valid mint eip712 signature, error: %v, valid:%v", err, valid)
		}
	})

	t.Run("meta tx mint", func(t1 *testing.T) {
		resp, err := client.Mint(contract, []string{"uri1"}, []common.Address{client.Address()}, mintSig)
		if err != nil {
			t1.Fatalf("failed to make metatx mint, error: %v", err)
		}
		if !resp.Success {
			t1.Fatalf("failed to make metatx mint, error: tx failed")
		}
		t1.Logf("mint meta tx | cost=%v | block number=%v", resp.Cost, resp.BlockNumber)
	})

	t.Run("meta tx mint fails for replay/invalid tx", func(t1 *testing.T) {
		_, err := client.Mint(contract, []string{"uri1"}, []common.Address{client.Address()}, mintSig)
		if err == nil {
			t1.Fatalf("failed catch replay make metatx mint")
		}
		if !strings.Contains(err.Error(), "tx signature mismatch") {
			t1.Fatalf("metatx mint failed for unrelated reason, error: %v", err)
		}
	})

	t.Run("verify user tokens", func(t1 *testing.T) {
		tokens, err := client.UserTokens(contract, client.Address())
		if err != nil {
			t1.Fatalf("failed to get tokens, error:%v", err)
		}
		expected := []string{"uri0", "uri1"}
		for id, uri := range tokens {
			if int(id) >= len(expected) {
				t1.Fatalf("token count mimatch, expected max id=%v got id=%v", len(expected)-1, id)
			}
			if expected[id] != uri {
				t1.Fatalf("metadata mismatch for id=%v, expected=%v got=%v", id, expected[id], uri)
			}
		}
	})
	t.Run("burn", func(t1 *testing.T) {
		resp, err := client.Burn(contract, []*big.Int{big.NewInt(0)})
		if err != nil {
			t1.Fatalf("failed to make burn, error: %v", err)
		}
		if !resp.Success {
			t1.Fatalf("failed to make burn, error: tx failed")
		}
		burnTx = resp.TxHash
		t1.Logf("burn tx | cost=%v | block number=%v", resp.Cost, resp.BlockNumber)
	})

	t.Run("burn logs", func(t1 *testing.T) {
		expectedData := BurnLogDetails{
			Ids:       []*big.Int{big.NewInt(0)},
			BurnNonce: big.NewInt(0),
			User:      client.Address(),
			Success:   true,
		}
		resp, err := client.BurnLog(contract, burnTx)
		if err != nil {
			t1.Fatalf("failed to get burn log, error: %v", err)
		}
		if resp.Success != expectedData.Success {
			t1.Fatalf("burn log fetched doesn't match expected value, expected=%v got=%v", expectedData.Success, resp.Success)
		}
		if resp.BurnNonce.Cmp(expectedData.BurnNonce) != 0 {
			t1.Fatalf("burn log fetched doesn't match expected value, expected=%v got=%v", expectedData.BurnNonce, resp.BurnNonce)
		}
		if resp.User != expectedData.User {
			t1.Fatalf("burn log fetched doesn't match expected value, expected=%v got=%v", expectedData.User, resp.User)
		}
		if len(resp.Ids) != len(expectedData.Ids) {
			t1.Fatalf("burn log fetched doesn't match expected value, expected len=%v got=%v", len(expectedData.Ids), len(resp.Ids))
		}
		for i := 0; i < len(resp.Ids); i++ {
			if resp.Ids[i].Cmp(expectedData.Ids[i]) != 0 {
				t1.Fatalf("burn log fetched doesn't match expected id value, expected=%v got=%v", expectedData.Ids[i], resp.Ids[i])
			}
		}
	})

	t.Run("check mint nonce", func(t1 *testing.T) {
		resp, err := signerClient.MintNonce(contract)
		if err != nil {
			t1.Fatalf("failed to get mint nonce, error: %v", err)
		}
		if resp.Cmp(big.NewInt(2)) != 0 {
			t1.Fatalf("invalid mint nonce, expected=%v got=%v", 2, resp.Int64())
		}
	})

	t.Run("check burn nonce", func(t1 *testing.T) {
		burnNonce, err = client.BurnNonce(contract)
		if err != nil {
			t1.Fatalf("failed to get burn nonce, error: %v", err)
		}
		if burnNonce.Cmp(big.NewInt(1)) != 0 {
			t1.Fatalf("invalid burn nonce, expected=%v got=%v", 1, burnNonce.Int64())
		}
	})

	t.Run("generate and verify burn signature", func(t1 *testing.T) {
		hash, err := client.GenerateBurnEIP712Hash(contract, burnNonce)
		if err != nil {
			t1.Fatalf("failed to get burn nonce, error: %v", err)
		}
		burnSig, err := signerClient.Sign(hash)
		if err != nil {
			t1.Fatalf("failed to sign burn eip712 hash, error: %v", err)
		}
		if valid, err := signerClient.VerifyEIP712Signature(hash, hex.EncodeToString(burnSig), signerClient.Address().String()); err != nil || !valid {
			t1.Fatalf("failed to generate valid burn eip712 signature, error: %v, valid:%v", err, valid)
		}
	})
}
