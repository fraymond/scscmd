# Some Customized Simple Go Command

This program will take a keystore file, a from address (must be in the given keystore) and to address, some ethereum amount and it will unlock the account and sign the transaction. If a passphrase is not provided, the user will be given 3 chances to put in the password. 

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

## Build SCS Command
```bash
go build scscmd.go
```

## Use SCS Command
```bash
./scscmd "/Users/[user]/Library/Ethereum/keystore" "94329ffd9c8c6651ed0250569700e823ad6ebcbd" "ca4a1cc346f6ed99a7a0335617a1c647244c767a" 500000000000000000000 "password"
```
