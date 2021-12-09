package main

import (
	"auction/fabric"
	"time"
)

type AuctionResult struct {
	HighestBidder []byte
}

func listenNewAuction() *fabric.Auction {
	lastID, err := assetCC.GetLastAuctionID()
	check(err)

	for {
		time.Sleep(1 * time.Second)
		auctionID, err := assetCC.GetLastAuctionID()
		check(err)
		if auctionID > lastID {
			auction, err := assetCC.GetAuction(auctionID)
			check(err)
			return auction
		}
	}
}

func listenAuctionEnding(auctionID int) bool {
	for {
		time.Sleep(1 * time.Second)
		auction, err := assetCC.GetAuction(auctionID)
		check(err)
		if auction.Status == "Ending" {
			return true
		}
	}
}
