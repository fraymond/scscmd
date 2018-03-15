package main

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/console"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/params"
)

var testTxPoolConfig core.TxPoolConfig

type testBlockChain struct {
	statedb       *state.StateDB
	gasLimit      uint64
	chainHeadFeed *event.Feed
}

func init() {
	testTxPoolConfig = core.DefaultTxPoolConfig
	testTxPoolConfig.Journal = ""
}

func (bc *testBlockChain) CurrentBlock() *types.Block {
	return types.NewBlock(&types.Header{
		GasLimit: bc.gasLimit,
	}, nil, nil, nil)
}

func (bc *testBlockChain) GetBlock(hash common.Hash, number uint64) *types.Block {
	return bc.CurrentBlock()
}

func (bc *testBlockChain) StateAt(common.Hash) (*state.StateDB, error) {
	return bc.statedb, nil
}

func (bc *testBlockChain) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return bc.chainHeadFeed.Subscribe(ch)
}

func transaction(nonce uint64, gaslimit uint64, key *ecdsa.PrivateKey) *types.Transaction {
	return pricedTransaction(nonce, gaslimit, big.NewInt(1), key)
}

func pricedTransaction(nonce uint64, gaslimit uint64, gasprice *big.Int, key *ecdsa.PrivateKey) *types.Transaction {
	tx, _ := types.SignTx(types.NewTransaction(nonce, common.Address{}, big.NewInt(100), gaslimit, gasprice, nil), types.HomesteadSigner{}, key)
	return tx
}

func setupTxPool(chainID *big.Int) (*core.TxPool, *testBlockChain) {
	diskdb, _ := ethdb.NewMemDatabase()
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(diskdb))
	blockchain := &testBlockChain{statedb, 200000000000000, new(event.Feed)}

	testChainConfig := params.TestChainConfig
	testChainConfig.ChainId = chainID

	//key, _ := crypto.GenerateKey()
	pool := core.NewTxPool(testTxPoolConfig, params.TestChainConfig, blockchain)

	return pool, blockchain
}

/* USAGE:
 * ./scscmd [keystore dir] [from account] [to account] [eth amount] [[passphrase]]
 */
func main() {

	// create a keystore from keystore file
	ks := keystore.NewKeyStore(os.Args[1], keystore.LightScryptN, keystore.LightScryptP)

	fromAddr := os.Args[2]
	toAddr := os.Args[3]
	amount := new(big.Int)
	amount, ok := amount.SetString(os.Args[4], 10)

	if !ok {
		fmt.Println("Invaid eth amount")
		return
	}

	var password string

	// if passphrase is not provided, the user is given 3 times to enter passphrase.
	if len(os.Args) == 6 {
		password = os.Args[5]
	} else {
		for trials := 0; trials < 3; trials++ {
			prompt := fmt.Sprintf("Unlocking account %s | Attempt %d/%d", fromAddr, trials+1, 3)
			password = getPassPhrase(prompt, false)
		}
	}

	// unlock account
	account := accounts.Account{Address: common.HexToAddress(fromAddr)}
	err := ks.Unlock(account, password)

	if err == nil {
		fmt.Println("---- from account ----")
		fmt.Println(account)
	} else {
		fmt.Println(err)
	}

	recipientAddr := common.HexToAddress(toAddr)

	// hardcoded. not sure what to do.
	chainID := big.NewInt(3)
	nonce := uint64(7)
	gasLimit := uint64(100000)
	gasPrice := big.NewInt(20) // 20 gwei

	// create a new transaction
	tx := types.NewTransaction(nonce, recipientAddr, amount, gasLimit, gasPrice, nil)

	// sign transaction
	signedTx, err := ks.SignTxWithPassphrase(account, password, tx, chainID)

	fmt.Println("---- signed transactions ----")
	fmt.Println(signedTx)

	var buff bytes.Buffer
	signedTx.EncodeRLP(&buff)

	fmt.Println("---- buff.Bytes ----")
	fmt.Printf("0x%x\n", buff.Bytes())

	var testSigData = make([]byte, 32)
	fmt.Println("---- test data ----")
	fmt.Println(testSigData)

	// add to pool
	pool, blockchain := setupTxPool(chainID)
	defer pool.Stop()

	blockchain.statedb.AddBalance(common.HexToAddress(fromAddr), big.NewInt(500000000))

	fmt.Print("Gas Price: ")
	fmt.Println(tx.GasPrice())
	fmt.Print("Gas: ")
	fmt.Println(tx.Gas())
	fmt.Print("Cost: ")
	fmt.Println(tx.Cost())

	fmt.Println(fromAddr)
	fmt.Println(blockchain.statedb.GetBalance(common.HexToAddress(fromAddr)))

	if err := pool.AddLocal(signedTx); err != nil {
		fmt.Println("---- error ----")
		fmt.Println(err)
	}

	fmt.Println(pool.Content())

	if err != nil {
		fmt.Println("---- err ----")
		fmt.Println(err)
		return
	}

	// sign hash
	signature, err := ks.SignHashWithPassphrase(account, password, testSigData)
	if err != nil {
		fmt.Println("---- err ----")
		fmt.Println(err)
		return
	}

	fmt.Println("---- signed hash ----")
	fmt.Println(signature)

}

// getPassPhrase retrieves the password associated with an account, either fetched
// from a list of preloaded passphrases, or requested interactively from the user.
func getPassPhrase(prompt string, confirmation bool) string {

	// Otherwise prompt the user for the password
	if prompt != "" {
		fmt.Println(prompt)
	}
	password, err := console.Stdin.PromptPassword("Passphrase: ")
	if err != nil {
		fmt.Printf("Failed to read passphrase: %v", err)
	}
	if confirmation {
		confirm, err := console.Stdin.PromptPassword("Repeat passphrase: ")
		if err != nil {
			fmt.Printf("Failed to read passphrase confirmation: %v", err)
		}
		if password != confirm {
			fmt.Printf("Passphrases do not match")
		}
	}
	return password
}
