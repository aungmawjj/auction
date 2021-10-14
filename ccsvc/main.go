package main

import (
	"auction/contract"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/auction_result", fetchAuctionResult)

	e.Logger.Fatal(e.Start(":9000"))
}

type AuctionResultRequest struct {
	Address []byte
}

type AuctionResult struct {
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
	eth, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		return err
	}
	cc, err := contract.NewAuction(common.BytesToAddress(req.Address), eth)
	if err != nil {
		return err
	}
	var result AuctionResult
	result.Ended, err = cc.Ended(&bind.CallOpts{})
	if err != nil {
		return err
	}
	highestBid, err := cc.HighestBid(&bind.CallOpts{})
	if err != nil {
		return err
	}
	result.HighestBid = highestBid.Int64()

	highestBidder, err := cc.HighestBidder(&bind.CallOpts{})
	if err != nil {
		return err
	}
	result.HighestBidder = highestBidder.Bytes()

	return c.JSON(http.StatusOK, result)
}
