package main

import (
	"auction/contract"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func printTxStatus(success bool) {
	if success {
		fmt.Println("Transaction successful")
	} else {
		fmt.Println("Transaction failed")
	}
}

func checkTx(client *ethclient.Client, hash common.Hash) (bool, error) {
	for {
		r, err := client.TransactionReceipt(context.Background(), hash)
		if err == ethereum.NotFound {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		if err != nil {
			return false, err
		}
		if r.Status == types.ReceiptStatusFailed {
			fmt.Printf("%+v\n", r)
		}
		return r.Status == types.ReceiptStatusSuccessful, nil
	}
}

func newAuctionSession(
	addr common.Address, eth *ethclient.Client, keyfile string,
) *contract.AuctionSession {
	cc, err := contract.NewAuction(addr, eth)
	check(err)
	return &contract.AuctionSession{
		Contract:     cc,
		TransactOpts: *newTransactor(keyfile),
	}
}

func newTransactor(keyfile string) *bind.TransactOpts {
	f, err := os.Open(keyfile)
	check(err)
	defer f.Close()
	auth, err := bind.NewTransactor(f, "password")
	check(err)
	auth.GasLimit = 1000000000000
	return auth
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
