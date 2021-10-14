package main

import (
	"auction/contract"
	"context"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func ExpEth() {
	eth := newEthClient()

	// init asset
	fmt.Println("\ninitializing asset")
	assetAddr, tx, _, err := contract.DeployAsset(newTransactor("keys/key0"), eth)
	check(err)
	success, err := checkTx(eth, tx.Hash())
	check(err)
	printTxStatus(success)
	if !success {
		panic("failed to deploy asset contract")
	}
	fmt.Println("asset address:", assetAddr.Hex())

	assetSession := newAssetSession(assetAddr, eth, "keys/key0")
	assetOwner, err := assetSession.Owner()
	check(err)
	fmt.Println("asset owner:", assetOwner.Hex())

	// init auction for the asset
	fmt.Println("\ninitializing auction for asset")
	auctionAddr, tx, _, err := contract.DeployAuction(newTransactor("keys/key0"), eth)
	check(err)
	success, err = checkTx(eth, tx.Hash())
	check(err)
	printTxStatus(success)
	if !success {
		panic("failed to deploy auction contract")
	}
	fmt.Printf("auction address: %s\n", auctionAddr.Hex())

	auctionSession := newAuctionSession(auctionAddr, eth, "keys/key1")
	beneficiary, err := auctionSession.Beneficiary()
	check(err)
	fmt.Println("auction beneficiary:", beneficiary.Hex())

	// bind auction address to asset
	fmt.Println("\nbinding auction to asset")
	tx, err = assetSession.StartAuction(auctionAddr)
	check(err)
	success, err = checkTx(eth, tx.Hash())
	check(err)
	printTxStatus(success)
	if !success {
		panic("failed to bind auction to asset")
	}
	assetAuction, err := assetSession.PendingAuction()
	check(err)
	fmt.Println("asset auction address:", assetAuction.Hex())

	// bid auction
	fmt.Println("\nbidding auction")
	auctionSession.TransactOpts.Value = big.NewInt(1000)
	tx, err = auctionSession.Bid()
	check(err)
	success, err = checkTx(eth, tx.Hash())
	check(err)
	printTxStatus(success)
	if !success {
		panic("failed to bid auction")
	}
	auctionSession.TransactOpts.Value = big.NewInt(0)

	highestBidder, err := auctionSession.HighestBidder()
	check(err)
	fmt.Println("highest bidder:", highestBidder.Hex())

	highestBid, err := auctionSession.HighestBid()
	check(err)
	fmt.Println("highest bid:", highestBid)

	// end auction
	fmt.Println("\nending auction")
	tx, err = auctionSession.EndAuction()
	check(err)
	success, err = checkTx(eth, tx.Hash())
	check(err)
	printTxStatus(success)
	if !success {
		panic("failed to end auction")
	}

	tx, err = assetSession.EndAuction(highestBidder)
	check(err)
	success, err = checkTx(eth, tx.Hash())
	check(err)
	printTxStatus(success)
	if !success {
		panic("failed to end auction for asset")
	}

	assetOwner, err = assetSession.Owner()
	check(err)
	fmt.Println("asset owner:", assetOwner.Hex())
}

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

func newEthClient() *ethclient.Client {
	c, err := ethclient.Dial(fmt.Sprintf("http://%s:8545", "localhost"))
	check(err)
	return c
}

func newAssetSession(
	addr common.Address, eth *ethclient.Client, keyfile string,
) *contract.AssetSession {
	cc, err := contract.NewAsset(addr, eth)
	check(err)
	return &contract.AssetSession{
		Contract:     cc,
		TransactOpts: *newTransactor(keyfile),
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
