package main

import (
	"auction/contract"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type AuctionResult struct {
	HighestBidder []byte
}

func listenAuctionEnd(auctionID []byte) (*AuctionResult, error) {
	auction, err := contract.NewAuction(common.BytesToAddress(auctionID), ethClient)
	if err != nil {
		return nil, err
	}
	opts := &bind.CallOpts{}

	time.Sleep(5 * time.Second)
	for {
		time.Sleep(1 * time.Second)
		ended, _ := auction.Ended(opts)
		if ended {
			break
		}
	}
	highestBidder, _ := auction.HighestBidder(opts)
	return &AuctionResult{
		HighestBidder: highestBidder.Bytes(),
	}, nil
}
