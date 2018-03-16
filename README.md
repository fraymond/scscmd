# Some Customized Simple Go Command

This program will use go program to send transactions to Ethereum network.

## scscmd.go

This program will take a keystore file, a from address (must be in the given keystore) and to address, some ethereum amount and it will unlock the account and sign the transaction. If a passphrase is not provided, the user will be given 3 chances to put in the password. 

It will also build a private test blockchain, and sent the signed transaction to the transaction pool of the test blockchain.

## ctrtcmd.go

This is a newer version. It uses ethclient to connect to a local Ethereum testnet using geth.ipc. It will test the following three types of transactions.

1. Send a balance.
2. Deploy a smart contract.
3. Transacting a smart contract.


This program uses [Go-Ethereum](https://github.com/ethereum/go-ethereum).

## Go-Ethereum Installation

### setup $GOPATH

```
export GOPATH=/Users/[user]/go
```

### go get
```bash
go get -u github.com/ethereum/go-ethereum
```

## Build and Use SCS Command
```bash
go build scscmd.go
./scscmd "/Users/[user]/Library/Ethereum/keystore" "94329ffd9c8c6651ed0250569700e823ad6ebcbd" "ca4a1cc346f6ed99a7a0335617a1c647244c767a" 500000000000000000000 "password"
```
## Prepare Smart Contract Coin.sol

Review JSON format of the contract
```bash
solc --optimize --combined-json abi,bin,interface Coin.sol
```
`
Create a JS file of the contract
```
echo "var testOutput=`solc --optimize --combined-json abi,bin,interface Coin.sol`" > coin.js
```

Open web3 console to load script
```
loadScript("coin.js")
```
Review contract abi and bin
```
testOutput.contracts["Coin.sol:Coin"]
```
Build the smart contract
```
var mintContract = eth.contract(JSON.parse(testOutput.contracts["Coin.sol:Coin"].abi));
```
Unlock Account
```
personal.unlockAccount(eth.accounts[0], "0000")
```
Deploy contract
```
var mint = mintContract.new({ from: eth.accounts[0], data: "0x" + testOutput.contracts["Coin.sol:Coin"].bin, gas: 4700000},
  function (e, contract) {
    console.log(e, contract);
    if (typeof contract.address !== 'undefined') {
         console.log('Contract mined! address: ' + contract.address + ' transactionHash: ' + contract.transactionHash);
    }
  }
);
```
Start mining to mine the contract
```
miner.start()
miner.stop()
```

## Run Ctrt Command
```bash
go run ctrtcmd.go
```