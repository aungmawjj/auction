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
var quorumClient *ethclient.Client
var assetCC *fabric.AssetCC

func Scanerio_1() {
	var err error
	ethClient, err = ethclient.Dial(fmt.Sprintf("http://%s:8545", "localhost"))
	check(err)

	quorumClient, err = ethclient.Dial(fmt.Sprintf("http://%s:8546", "localhost"))
	check(err)

	assetCC = fabric.NewAssetCC()

	auctionWithQuorum()
	auctionWithEthereum()
}

func auctionWithEthereum() {
	var err error
	fmt.Printf("\nAuction with Quorum\n\n")

	fmt.Println("[fabric] Adding asset")
	asset := addAsset("asset1")

	fmt.Println("[ccsvc] Creating auction for asset")
	createAuction([]byte(assetCC.GetCCID()), asset.ID, "ethereum")

	asset, err = assetCC.GetAsset(asset.ID)
	check(err)

	auctionID := asset.PendingAuction.ID
	auctionAddr := common.BytesToAddress(auctionID)
	fmt.Println("auction ID: ", auctionAddr.Hex())

	fmt.Println("\n[ethereum] bidding auction")
	bidAuction(auctionAddr, ethClient)

	fmt.Println("\n[ethereum] ending auction")
	endAuction(auctionAddr, ethClient)

	time.Sleep(15 * time.Second)

	asset, err = assetCC.GetAsset(asset.ID)
	check(err)
	fmt.Println("Asset Owner:", common.BytesToAddress(asset.Owner).Hex())
}

func auctionWithQuorum() {
	var err error
	fmt.Printf("\nAuction with Quorum\n\n")

	fmt.Println("[fabric] Adding asset")
	asset := addAsset("asset2")

	fmt.Println("[ccsvc] Creating auction for asset, platform: quorum")
	createAuction([]byte(assetCC.GetCCID()), asset.ID, "quorum")

	asset, err = assetCC.GetAsset(asset.ID)
	check(err)

	auctionID := asset.PendingAuction.ID
	auctionAddr := common.BytesToAddress(auctionID)
	fmt.Println("auction ID: ", auctionAddr.Hex())

	fmt.Println("\n[quorum] bidding auction")
	bidAuction(auctionAddr, quorumClient)

	fmt.Println("\n[quorum] ending auction")
	endAuction(auctionAddr, quorumClient)

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

func bidAuction(addr common.Address, client *ethclient.Client) {
	auctionSession := newAuctionSession(addr, client, "keys/key1")
	auctionSession.TransactOpts.Value = big.NewInt(1000)
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

func endAuction(addr common.Address, client *ethclient.Client) {
	auctionSession := newAuctionSession(addr, client, "keys/key0")
	tx, err := auctionSession.EndAuction()
	check(err)
	success, err := checkTx(client, tx.Hash())
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
