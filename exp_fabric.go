package main

import (
	"auction/contract"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"
)

func ExpFabric() {

	fabtic := NewFabricClient("http://localhost:7050")
	fabtic.ChaincodePath = "github.com/aungmawjj/crosschain_cc"

	assetCC := &AssetCC{
		fabric: fabtic,
	}

	ccid, err := assetCC.Deploy()
	check(err)
	fabtic.ChaincodeID = ccid

	time.Sleep(5 * time.Second)
	asset := Asset{
		ID:    sha256Sum("asset1"),
		Owner: newTransactor("keys/key0").From.Bytes(),
	}
	fmt.Println("[fabric] Adding asset, owner: ", hex.EncodeToString(asset.Owner))
	txID, err := assetCC.AddAsset(asset)
	check(err)
	check(fabtic.WaitTx(txID))

	asset, err = assetCC.GetAsset(asset.ID)
	check(err)
	fmt.Println("Asset added, owner: ", hex.EncodeToString(asset.Owner))

	eth := newEthClient()

	// init auction for the asset
	fmt.Println("\n[ethereum] initializing auction for asset")
	auctionAddr, tx, _, err := contract.DeployAuction(newTransactor("keys/key0"), eth)
	check(err)
	success, err := checkTx(eth, tx.Hash())
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
	fmt.Println("\n[fabric] binding auction to asset")
	txID, err = assetCC.SetAuction(SetAuctionArgs{
		AssetID:   asset.ID,
		AuctionID: auctionAddr.Bytes(),
	})
	check(err)
	check(fabtic.WaitTx(txID))

	// bid auction
	fmt.Println("\n[ethereum] bidding auction")
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
	fmt.Println("\n[ethereum] ending auction")
	tx, err = auctionSession.EndAuction()
	check(err)
	success, err = checkTx(eth, tx.Hash())
	check(err)
	printTxStatus(success)
	if !success {
		panic("failed to end auction")
	}

	// end auction on fabric
	fmt.Println("\n[fabric] ending auction")
	txID, err = assetCC.EndAuction(asset.ID)
	check(err)
	check(fabtic.WaitTx(txID))

	asset, err = assetCC.GetAsset(asset.ID)
	check(err)
	fmt.Printf("Asset Owner: %s", hex.EncodeToString(asset.Owner))
}

func sha256Sum(data string) []byte {
	sum := sha256.Sum256([]byte(data))
	return sum[:]
}
