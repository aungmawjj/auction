package main

import (
	"encoding/hex"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var auctionResults map[string]*AuctionResult

func main() {

	auctionResults = make(map[string]*AuctionResult)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/auction_result", fetchAuctionResult)
	e.POST("/on_receive", handleOnReceiveAuctionResult)

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
}

func fetchAuctionResult(c echo.Context) error {
	var req AuctionResultRequest
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	result, ok := auctionResults[string(req.Address)]
	if !ok {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, result)
}

// func fetchAuctionResult(c echo.Context) error {
// 	var req AuctionResultRequest
// 	err := c.Bind(&req)
// 	if err != nil {
// 		return err
// 	}
// 	eth, err := ethclient.Dial("http://localhost:8545")
// 	if err != nil {
// 		return err
// 	}
// 	cc, err := contract.NewAuction(common.BytesToAddress(req.Address), eth)
// 	if err != nil {
// 		return err
// 	}
// 	var result AuctionResult
// 	result.Ended, err = cc.Ended(&bind.CallOpts{})
// 	if err != nil {
// 		return err
// 	}
// 	highestBid, err := cc.HighestBid(&bind.CallOpts{})
// 	if err != nil {
// 		return err
// 	}
// 	result.HighestBid = highestBid.Int64()

// 	highestBidder, err := cc.HighestBidder(&bind.CallOpts{})
// 	if err != nil {
// 		return err
// 	}
// 	result.HighestBidder = highestBidder.Bytes()

// 	return c.JSON(http.StatusOK, result)
// }

func onReceiveAuctionResult(result *AuctionResult) {
	log.Printf("Received auction result, %s\n", hex.EncodeToString(result.Address))
	auctionResults[string(result.Address)] = result
}
