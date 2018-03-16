package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	// SendBalance()
	// DeployContract()
	CallContract()
}

// SendBalance is used to send a balance as transaction
func SendBalance() {

	// Dial IPC
	client, err := ethclient.Dial("/Users/rfu/geth-test/privatebc/geth.ipc")

	// Dial RPC
	//client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}
	fmt.Println(client)

	// 1. Send a Balance Transaction.
	d := time.Now().Add(10000 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	fromAddr := "0x18833df6ba69b4d50acc744e8294d128ed8db1f1"
	toAddr := "0x882dbeb3de07f01df95e14e9db16d834a8ceea8f"
	password := "0000"

	account := accounts.Account{Address: common.HexToAddress(fromAddr)}

	ks := keystore.NewKeyStore("/Users/rfu/geth-test/privatebc/keystore", keystore.LightScryptN, keystore.LightScryptP)
	_, err = ks.Find(account)
	if err != nil {
		log.Fatalf("Can't find account: %v", err)
	}
	err = ks.Unlock(account, password)

	if err != nil {
		log.Fatalf("Can't unlock account: %v", err)
	}

	if err != nil {
		log.Fatalf("unlock failed: %v", err)
	}

	recipientAddr := common.HexToAddress(toAddr)

	nonce, _ := client.NonceAt(ctx, common.HexToAddress(fromAddr), nil)

	// hardcoded. not sure what to do.
	gasLimit := uint64(100000)
	gasPrice := big.NewInt(20) // 20 gwei
	amount := big.NewInt(100)
	chainID := big.NewInt(4321) // my private blockchain used network id 4321

	// create a new transaction
	tx := types.NewTransaction(nonce, recipientAddr, amount, gasLimit, gasPrice, nil)

	// sign transaction
	signedTx, err := ks.SignTxWithPassphrase(account, password, tx, chainID)

	fmt.Println("---- signed transactions ----")
	fmt.Println(signedTx)

	err = client.SendTransaction(ctx, signedTx)

	if err != nil {
		fmt.Println(err)
	}
}

// DeployContract is deploying a contract
func DeployContract() {

	// Dial IPC
	client, err := ethclient.Dial("/Users/rfu/geth-test/privatebc/geth.ipc")

	// Dial RPC
	//client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}
	fmt.Println(client)

	// 1. Send a Balance Transaction.
	d := time.Now().Add(10000 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	fromAddr := "0x18833df6ba69b4d50acc744e8294d128ed8db1f1"
	//toAddr := "0x882dbeb3de07f01df95e14e9db16d834a8ceea8f"
	password := "0000"

	account := accounts.Account{Address: common.HexToAddress(fromAddr)}

	ks := keystore.NewKeyStore("/Users/rfu/geth-test/privatebc/keystore", keystore.LightScryptN, keystore.LightScryptP)
	_, err = ks.Find(account)
	if err != nil {
		log.Fatalf("Can't find account: %v", err)
	}
	err = ks.Unlock(account, password)

	if err != nil {
		log.Fatalf("Can't unlock account: %v", err)
	}

	if err != nil {
		log.Fatalf("unlock failed: %v", err)
	}

	//recipientAddr := common.HexToAddress(toAddr)

	nonce, _ := client.NonceAt(ctx, common.HexToAddress(fromAddr), nil)

	// hardcoded. not sure what to do.
	gasLimit := uint64(2000000)
	gasPrice := big.NewInt(30)  // 20 gwei
	chainID := big.NewInt(4321) // my private blockchain used network id 4321

	// 2. Call a Contract
	coinAbi := "[{\"constant\":true,\"inputs\":[],\"name\":\"minter\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"receivers\",\"type\":\"address[]\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"send\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Sent\",\"type\":\"event\"}]"
	coinBin := `0x6060604052341561000f57600080fd5b60008054600160a060020a033316600160a060020a03199091161790556102ca8061003b6000396000f3006060604052600436106100615763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166307546172811461006657806327e235e31461009557806340c10f19146100c6578063efd5a170146100ea575b600080fd5b341561007157600080fd5b61007961013b565b604051600160a060020a03909116815260200160405180910390f35b34156100a057600080fd5b6100b4600160a060020a036004351661014a565b60405190815260200160405180910390f35b34156100d157600080fd5b6100e8600160a060020a036004351660243561015c565b005b34156100f557600080fd5b6100e86004602481358181019083013580602081810201604051908101604052809392919081815260200183836020028082843750949650509335935061019a92505050565b600054600160a060020a031681565b60016020526000908152604090205481565b60005433600160a060020a0390811691161461017757610196565b600160a060020a03821660009081526001602052604090208054820190555b5050565b60008251600160a060020a0333166000908152600160205260409020549083029010156101c657610299565b5060005b825181101561029957600160a060020a0333166000908152600160208190526040822080548590039055839185848151811061020257fe5b90602001906020020151600160a060020a031681526020810191909152604001600020805490910190557f3990db2d31862302a685e8086b5755072a6e2b5b780af1ee81ece35ee3cd33453384838151811061025a57fe5b9060200190602002015184604051600160a060020a039384168152919092166020820152604080820192909252606001905180910390a16001016101ca565b5050505600a165627a7a723058206d07780b0454dad22a4068afdaf66dfc31075ee58e14699ae4486e363481a8c50029`

	byteCode := common.Hex2Bytes(coinBin[2:])
	cAbi, _ := abi.JSON(strings.NewReader(coinAbi))
	input, _ := cAbi.Pack("")

	byteCode = append(byteCode, input...)

	tx := types.NewContractCreation(nonce, nil, gasLimit, gasPrice, byteCode)

	// sign transaction
	signedTx, err := ks.SignTxWithPassphrase(account, password, tx, chainID)

	fmt.Println("---- signed transactions ----")
	fmt.Println(signedTx)

	err = client.SendTransaction(ctx, signedTx)

	if err != nil {
		fmt.Println(err)
	}
}

// CallContract is calling a function from a contract
func CallContract() {

	// Dial IPC
	client, err := ethclient.Dial("/Users/rfu/geth-test/privatebc/geth.ipc")

	// Dial RPC
	//client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum client: %v", err)
	}
	fmt.Println(client)

	d := time.Now().Add(10000 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	fromAddr := "0x18833df6ba69b4d50acc744e8294d128ed8db1f1"

	receivers := [2]string{"0x18833df6ba69b4d50acc744e8294d128ed8db1f1", "0x882dbeb3de07f01df95e14e9db16d834a8ceea8f"}
	password := "0000"

	account := accounts.Account{Address: common.HexToAddress(fromAddr)}

	ks := keystore.NewKeyStore("/Users/rfu/geth-test/privatebc/keystore", keystore.LightScryptN, keystore.LightScryptP)
	_, err = ks.Find(account)
	if err != nil {
		log.Fatalf("Can't find account: %v", err)
	}
	err = ks.Unlock(account, password)

	if err != nil {
		log.Fatalf("Can't unlock account: %v", err)
	}

	if err != nil {
		log.Fatalf("unlock failed: %v", err)
	}

	//recipientAddr := common.HexToAddress(toAddr)

	nonce, _ := client.NonceAt(ctx, common.HexToAddress(fromAddr), nil)

	// hardcoded. not sure what to do.
	gasLimit := uint64(20000000)
	gasPrice := big.NewInt(30)  // 20 gwei
	chainID := big.NewInt(4321) // my private blockchain used network id 4321

	// 2. Call a Contract
	coinAbi := "[{\"constant\":true,\"inputs\":[],\"name\":\"minter\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"receiver\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"receivers\",\"type\":\"address[]\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"send\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Sent\",\"type\":\"event\"}]"
	// coinBin := `0x6060604052341561000f57600080fd5b60008054600160a060020a033316600160a060020a03199091161790556102ca8061003b6000396000f3006060604052600436106100615763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166307546172811461006657806327e235e31461009557806340c10f19146100c6578063efd5a170146100ea575b600080fd5b341561007157600080fd5b61007961013b565b604051600160a060020a03909116815260200160405180910390f35b34156100a057600080fd5b6100b4600160a060020a036004351661014a565b60405190815260200160405180910390f35b34156100d157600080fd5b6100e8600160a060020a036004351660243561015c565b005b34156100f557600080fd5b6100e86004602481358181019083013580602081810201604051908101604052809392919081815260200183836020028082843750949650509335935061019a92505050565b600054600160a060020a031681565b60016020526000908152604090205481565b60005433600160a060020a0390811691161461017757610196565b600160a060020a03821660009081526001602052604090208054820190555b5050565b60008251600160a060020a0333166000908152600160205260409020549083029010156101c657610299565b5060005b825181101561029957600160a060020a0333166000908152600160208190526040822080548590039055839185848151811061020257fe5b90602001906020020151600160a060020a031681526020810191909152604001600020805490910190557f3990db2d31862302a685e8086b5755072a6e2b5b780af1ee81ece35ee3cd33453384838151811061025a57fe5b9060200190602002015184604051600160a060020a039384168152919092166020820152604080820192909252606001905180910390a16001016101ca565b5050505600a165627a7a723058206d07780b0454dad22a4068afdaf66dfc31075ee58e14699ae4486e363481a8c50029`

	contractAddress := "0x5E7908fd90288a7ec0d0cbf9eC143377A9f4D574"
	cAbi, _ := abi.JSON(strings.NewReader(coinAbi))
	input, _ := cAbi.Pack("send", receivers, 100)

	// from address is contract address, amount is nil
	tx := types.NewTransaction(nonce, common.HexToAddress(contractAddress), nil, gasLimit, gasPrice, input)

	// sign transaction
	signedTx, err := ks.SignTxWithPassphrase(account, password, tx, chainID)

	fmt.Println("---- signed transactions ----")
	fmt.Println(signedTx)

	err = client.SendTransaction(ctx, signedTx)

	if err != nil {
		fmt.Println(err)
	}
}
