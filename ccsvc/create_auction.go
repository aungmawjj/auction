package main

import (
	"auction/contract"
	"auction/events"
	"auction/fabric"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Shopify/sarama"
	"github.com/labstack/echo/v4"
)

type CreateAuctionRequest struct {
	AssetCC  []byte
	AssetID  []byte
	Platform string
}

func handleCreateAuction(c echo.Context) error {
	var req CreateAuctionRequest
	err := c.Bind(&req)
	if err != nil {
		return err
	}
	go createAuction(req)
	return c.NoContent(http.StatusOK)
}

func createAuction(req CreateAuctionRequest) {
	var err error
	log.Println("Creating auction for asset")

	addr, _, _, err := contract.DeployAuction(ethTransactor, ethClient)
	if err != nil {
		log.Printf("failed to deploy auction %+v", err)
		return
	}
	log.Printf("Deployed auction on ethereum: %s", addr.Hex())

	fabtic := fabric.NewFabricClient("http://localhost:7050")
	fabtic.ChaincodeID = string(req.AssetCC)
	assetCC := fabric.NewAssetCC(fabtic)

	args := fabric.BindAuctionArgs{
		AssetID: req.AssetID,
		Auction: fabric.Auction{
			ID: addr.Bytes(),
		},
	}
	assetCC.BindAuction(args)

	log.Printf("Bind auction on fabric")

	message := &sarama.ProducerMessage{Topic: events.TopicOnBindAuction, Partition: 0}
	value, _ := json.Marshal(events.OnBindAuction{
		AssetCC:   req.AssetCC,
		AssetID:   req.AssetID,
		AuctionID: addr.Bytes(),
	})
	message.Value = sarama.ByteEncoder(value)

	_, _, err = kafkaProducer.SendMessage(message)
	if err != nil {
		log.Printf("failed to send kafka message %+v", err)
		return
	}
	log.Printf("Published event, OnBindAuction")
}
