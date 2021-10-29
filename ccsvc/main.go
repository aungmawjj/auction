package main

import (
	"auction/fabric"
	"encoding/hex"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()
	e.Use(middleware.Recover())

	e.POST("/on_receive", handleOnReceiveAuctionResult)

	Subscribe()
	e.Logger.Fatal(e.Start(":9000"))
}

type AuctionResultRequest struct {
	Address []byte
}

type AuctionResult struct {
	Address       []byte
	Ended         bool
	HighestBid    int64
	HighestBidder []byte
	FabricCCID    string
}

func onReceiveAuctionResult(result *AuctionResult) {
	log.Printf("Received auction result, %s\n", hex.EncodeToString(result.Address))

	fabtic := fabric.NewFabricClient("http://localhost:7050")
	fabtic.ChaincodePath = "github.com/aungmawjj/crosschain_cc"
	fabtic.ChaincodeID = result.FabricCCID
	assetCC := fabric.NewAssetCC(fabtic)

	auction := fabric.Auction{
		ID:            result.Address,
		Ended:         result.Ended,
		HighestBid:    result.HighestBid,
		HighestBidder: result.HighestBidder,
	}

	_, err := assetCC.UpdateAuction(auction)
	if err != nil {
		log.Printf("failed to update auction on fabric %+v", err)
	}
}
