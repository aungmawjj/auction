package main

import (
	"auction/fabric"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type CreateAuctionRequest struct {
	AssetCC  []byte
	AssetID  []byte
	Platform string
}

var ethClient *ethclient.Client
var quorumClient *ethclient.Client
var assetCC *fabric.AssetCC

func ThreeChainAuction() {
	var err error
	ethClient, err = ethclient.Dial(fmt.Sprintf("http://%s:8545", "localhost"))
	check(err)

	quorumClient, err = ethclient.Dial(fmt.Sprintf("http://%s:8546", "localhost"))
	check(err)

	assetCC = fabric.NewAssetCC()

	var asset *fabric.Asset
	var auction *fabric.Auction

	fmt.Println("[fabric] Adding asset")
	asset = addAsset("asset1")

	fmt.Println("[fabric] Start auction")
	auction = startAuction(asset.ID, []string{"ethereum", "quorum"})

	fmt.Println("[ethereum] Bid auction")
	// bidAuction(ethClient, auction.CrossAuctionIDs[0], "keys/key1", 500)

	fmt.Println("[quorum] Bid auction")
	// bidAuction(quorumClient, auction.CrossAuctionIDs[1], "keys/key2", 1000)

	fmt.Println("[fabric] End auction")
	endAuction(auction)
}

func addAsset(id string) *fabric.Asset {
	_, err := assetCC.AddAsset(id, newTransactor("keys/key0").From.Hex())
	check(err)
	time.Sleep(3 * time.Second)
	asset, err := assetCC.GetAsset(id)
	check(err)
	fmt.Println("Asset added, owner: ", asset.Owner)
	return asset
}

func startAuction(assetID string, platforms []string) *fabric.Auction {
	args := fabric.StartAuctionArgs{
		AssetID:   assetID,
		Platforms: platforms,
	}
	_, err := assetCC.StartAuction(args)
	check(err)
	time.Sleep(3 * time.Second)
	fmt.Println("Started auction for asset")

	auctionID, err := assetCC.GetLastAuctionID()
	check(err)
	fmt.Println("AuctionID: ", auctionID)

	// mock
	mockBindAuction(auctionID)

	for {
		time.Sleep(1 * time.Second)
		auction, err := assetCC.GetAuction(auctionID)
		check(err)
		if auction.Status == "Bind" {
			fmt.Println("Cross-chain auctions bind successful")
			return auction
		}
	}
}

func endAuction(auction *fabric.Auction) {
	_, err := assetCC.SetAuctionEnding(auction.AssetID)
	check(err)

	// mock
	mockEndAuction(auction.ID)

	for {
		time.Sleep(1 * time.Second)
		auction, err = assetCC.GetAuction(auction.ID)
		check(err)
		if auction.Status == "Ended" {

			fmt.Println("Auction Ended")
			fmt.Println("Highest Bidder: ", auction.HighestBidder)
			fmt.Println("Highest Bid: ", auction.HighestBid)
			fmt.Println("Highest Bid Platform: ", auction.HighestBidPlatform)

			asset, err := assetCC.GetAsset(auction.AssetID)
			check(err)
			fmt.Println("Asset Owner: ", asset.Owner)

			break
		}
	}
}

func bidAuction(client *ethclient.Client, addrHex, keyfile string, value int64) {
	addr := common.HexToAddress(addrHex)
	auctionSession := newAuctionSession(addr, client, keyfile)
	auctionSession.TransactOpts.Value = big.NewInt(value)
	tx, err := auctionSession.Bid()
	check(err)
	success, err := checkTx(client, tx.Hash())
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
}

func mockBindAuction(auctionID int) {
	args := fabric.BindAuctionArgs{
		AuctionID:       auctionID,
		CrossAuctionIDs: []string{"1", "2"},
	}
	_, err := assetCC.BindAuction(args)
	check(err)
}

func mockEndAuction(auctionID int) {
	args := fabric.EndAuctionArgs{
		AuctionID:      auctionID,
		HighestBids:    []int{500, 1000},
		HighestBidders: []string{"bidder1", "bidder2"},
	}
	_, err := assetCC.EndAuction(args)
	check(err)
}
