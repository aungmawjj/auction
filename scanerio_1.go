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
	fabtic := fabric.NewFabricClient("http://localhost:7050")
	fabtic.ChaincodePath = "github.com/aungmawjj/crosschain_cc"
	assetCC = fabric.NewAssetCC(fabtic)

	fmt.Println("[fabric] deploying asset chaincode")
	deployAssetCC()

	fmt.Println("[fabric] Adding asset")
	asset := addAsset()

	fmt.Println("[ccsvc] Creating auction for asset")
	createAuction([]byte(assetCC.GetCCID()), asset.ID)

	asset, err := assetCC.GetAsset(asset.ID)
	check(err)

	auctionID := asset.PendingAuction.ID
	auctionAddr := common.BytesToAddress(auctionID)

	fmt.Println("\n[ethereum] bidding auction")
	bidAuction(auctionAddr)

	fmt.Println("\n[ethereum] ending auction")
	endAuction(auctionAddr)

	time.Sleep(10 * time.Second)

	asset, err = assetCC.GetAsset(asset.ID)
	check(err)
	fmt.Println("Asset Owner:", hex.EncodeToString(asset.Owner))

}

func deployAssetCC() {
	ccid, err := assetCC.Deploy()
	check(err)
	fmt.Println("chaincode id:", ccid)
	assetCC.SetChaincodeID(ccid)
	time.Sleep(15 * time.Second)
}

func addAsset() fabric.Asset {
	asset := fabric.Asset{
		ID:    sha256Sum("asset1"),
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

func createAuction(assetCC, assetID []byte) {
	buf := bytes.NewBuffer(nil)
	json.NewEncoder(buf).Encode(CreateAuctionRequest{
		AssetCC: []byte(assetCC),
		AssetID: assetID,
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
