package main

import (
	"auction/contract"
	"auction/fabric"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"
)

func ExpFabric() {
	assetCC := fabric.NewAssetCC()

	var err error
	time.Sleep(20 * time.Second)
	asset := fabric.Asset{
		ID:    sha256Sum("asset1"),
		Owner: newTransactor("keys/key0").From.Bytes(),
	}
	fmt.Println("[fabric] Adding asset, owner: ", hex.EncodeToString(asset.Owner))
	_, err = assetCC.AddAsset(asset)
	check(err)
	time.Sleep(5 * time.Second)

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
	_, err = assetCC.BindAuction(fabric.BindAuctionArgs{
		AssetID: asset.ID,
		Auction: fabric.Auction{
			ID: auctionAddr.Bytes(),
		},
	})
	check(err)
	time.Sleep(5 * time.Second)

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

	err = fetchAuctionResult(auctionAddr.Bytes())
	check(err)

	// end auction on fabric
	fmt.Println("\n[fabric] ending auction")
	_, err = assetCC.EndAuction(fabric.EndAuctionArgs{})
	check(err)
	time.Sleep(5 * time.Second)

	asset, err = assetCC.GetAsset(asset.ID)
	check(err)
	fmt.Println("Asset Owner:", hex.EncodeToString(asset.Owner))
}
