package main

import (
	"auction/events"
	"auction/fabric"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Shopify/sarama"

	kb "github.com/philipjkim/kafka-brokers-go"
)

var assetCC *fabric.AssetCC
var kafkaProducer sarama.SyncProducer

func main() {
	fmt.Println("Fabric listener")
	var err error
	assetCC = fabric.NewAssetCC()

	zkNodes := []string{"localhost:2181"}
	kbConn, err := kb.NewConn(zkNodes)
	check(err)
	brokerList, _, err := kbConn.GetW()
	check(err)

	setupKafkaProducer(brokerList)
	go runNewAuctionListener()
	select {}
}

func setupKafkaProducer(brokerList []string) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewManualPartitioner

	var err error
	kafkaProducer, err = sarama.NewSyncProducer(brokerList, config)
	check(err)
}

func runNewAuctionListener() {
	for {
		auction := listenNewAuction()
		go onStartedAuction(auction)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func onStartedAuction(auction *fabric.Auction) {
	fmt.Println("Started new auction for asset: ", auction.AssetID)
	publishOnStartAuctionEvent(auction)

	listenAuctionEnding(auction.ID)
	fmt.Println("Ending auction, ID: ", auction.ID)
	publishOnEndAuctionEvent(auction)
}

func publishOnEndAuctionEvent(auction *fabric.Auction) {
	message := &sarama.ProducerMessage{Topic: events.TopicOnEndAuction, Partition: 0}
	value, _ := json.Marshal(events.AuctionEventPayload{
		AuctionID: auction.ID,
	})
	message.Value = sarama.ByteEncoder(value)
	_, _, err := kafkaProducer.SendMessage(message)
	if err != nil {
		log.Printf("failed to send kafka message %+v", err)
	}
	check(err)
	fmt.Println("Published event: ", events.TopicOnEndAuction)
}

func publishOnStartAuctionEvent(auction *fabric.Auction) {
	message := &sarama.ProducerMessage{Topic: events.TopicOnStartAuction, Partition: 0}
	value, _ := json.Marshal(events.AuctionEventPayload{
		AuctionID: auction.ID,
	})
	message.Value = sarama.ByteEncoder(value)
	_, _, err := kafkaProducer.SendMessage(message)
	if err != nil {
		log.Printf("failed to send kafka message %+v", err)
	}
	check(err)

	fmt.Println("Published event: ", events.TopicOnStartAuction)
}
