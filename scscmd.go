package main

import (
	"bytes"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/console"
	"github.com/ethereum/go-ethereum/core/types"
)

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
	gasPrice := big.NewInt(20000000000) // 20 gwei

	// create a new transaction
	tx := types.NewTransaction(nonce, recipientAddr, amount, gasLimit, gasPrice, nil)

	// sign transaction
	signedTx, err := ks.SignTxWithPassphrase(account, password, tx, chainID)

	if err != nil {
		fmt.Println("---- err ----")
		fmt.Println(err)
		return
	}

	fmt.Println("---- signed transactions ----")
	fmt.Println(signedTx)

	var buff bytes.Buffer
	signedTx.EncodeRLP(&buff)

	fmt.Println("---- buff.Bytes ----")
	fmt.Printf("0x%x\n", buff.Bytes())
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
