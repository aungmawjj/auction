package main

import (
	"auction/fabric"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
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
var assetCC *fabric.AssetCC

func Scanerio_1() {
	ethClient = newEthClient()
	assetCC = fabric.NewAssetCC()

	fmt.Println("[fabric] Adding asset")
	asset := addAsset("asset1")

	fmt.Println("[ccsvc] Creating auction for asset")
	createAuction([]byte(assetCC.GetCCID()), asset.ID, "ethereum")

	asset, err := assetCC.GetAsset(asset.ID)
	check(err)

	auctionID := asset.PendingAuction.ID
	auctionAddr := common.BytesToAddress(auctionID)
	fmt.Println("auction ID: ", auctionAddr.Hex())

	fmt.Println("\n[ethereum] bidding auction")
	bidAuction(auctionAddr)

	fmt.Println("\n[ethereum] ending auction")
	endAuction(auctionAddr)

	time.Sleep(15 * time.Second)

	asset, err = assetCC.GetAsset(asset.ID)
	check(err)
	fmt.Println("Asset Owner:", common.BytesToAddress(asset.Owner).Hex())

	fmt.Println("[fabric] Adding asset")
	asset = addAsset("asset2")

	fmt.Println("[ccsvc] Creating auction for asset, platform: quorum")
	createAuction([]byte(assetCC.GetCCID()), asset.ID, "quorum")

	asset, err = assetCC.GetAsset(asset.ID)
	check(err)

	auctionID = asset.PendingAuction.ID
	auctionAddr = common.BytesToAddress(auctionID)
	fmt.Println("auction ID: ", auctionAddr.Hex())

	fmt.Println("\n[quorum] bidding auction")
	bidAuction(auctionAddr)

	fmt.Println("\n[quorum] ending auction")
	endAuction(auctionAddr)

	time.Sleep(15 * time.Second)

	asset, err = assetCC.GetAsset(asset.ID)
	check(err)
	fmt.Println("Asset Owner:", common.BytesToAddress(asset.Owner).Hex())
}

func addAsset(id string) fabric.Asset {
	asset := fabric.Asset{
		ID:    sha256Sum(id),
		Owner: newTransactor("keys/key0").From.Bytes(),
	}
	_, err := assetCC.AddAsset(asset)
	check(err)
	time.Sleep(5 * time.Second)
	asset, err = assetCC.GetAsset(asset.ID)
	check(err)
	fmt.Println("Asset added, owner: ", hex.EncodeToString(asset.Owner))
	return asset
}

func createAuction(assetCC, assetID []byte, platform string) {
	buf := bytes.NewBuffer(nil)
	json.NewEncoder(buf).Encode(CreateAuctionRequest{
		AssetCC:  []byte(assetCC),
		AssetID:  assetID,
		Platform: platform,
	})

	resp, err := http.Post("http://localhost:9000/auction", "application/json", buf)
	check(err)
	resp.Body.Close()
	time.Sleep(10 * time.Second)
}

func bidAuction(addr common.Address) {
	auctionSession := newAuctionSession(addr, ethClient, "keys/key1")
	auctionSession.TransactOpts.Value = big.NewInt(1000)
	tx, err := auctionSession.Bid()
	check(err)
	success, err := checkTx(ethClient, tx.Hash())
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

func endAuction(addr common.Address) {
	auctionSession := newAuctionSession(addr, ethClient, "keys/key0")
	tx, err := auctionSession.EndAuction()
	check(err)
	success, err := checkTx(ethClient, tx.Hash())
	check(err)
	printTxStatus(success)
	if !success {
		panic("failed to end auction")
	}
}

func sha256Sum(data string) []byte {
	sum := sha256.Sum256([]byte(data))
	return sum[:]
}
