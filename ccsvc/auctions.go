package main

import (
	"auction/contract"
	"log"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func deployCrossChainAuctions() ([]string, error) {
	ethAddr, _, _, err := contract.DeployAuction(ethTransactor, ethClient)
	if err != nil {
		log.Printf("failed to deploy auction %+v", err)
		return nil, err
	}
	log.Printf("Deployed auction on ethereum: %s", ethAddr.Hex())

	quorumAddr, _, _, err := contract.DeployAuction(ethTransactor, quorumClient)
	if err != nil {
		log.Printf("failed to deploy auction %+v", err)
		return nil, err
	}
	log.Printf("Deployed auction on ethereum: %s", quorumAddr.Hex())

	return []string{ethAddr.Hex(), quorumAddr.Hex()}, nil
}

func getAuctionInfo(addr string, client *ethclient.Client) (int, string) {
	auction, err := contract.NewAuction(common.HexToAddress(addr), client)
	check(err)

	opts := &bind.CallOpts{}
	highestBid, err := auction.HighestBid(opts)
	check(err)

	highestBidder, err := auction.HighestBidder(opts)
	check(err)

	return int(highestBid.Int64()), highestBidder.Hex()
}
