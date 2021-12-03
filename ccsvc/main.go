package main

import (
	"auction/events"
	"auction/fabric"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	kb "github.com/philipjkim/kafka-brokers-go"
	"github.com/wvanbergen/kafka/consumergroup"
)

const kafkaConsumerGroup = "ccsvc"
const privKeyFile = "../keys/key0"

var ethClient *ethclient.Client
var quorumClient *ethclient.Client
var ethTransactor *bind.TransactOpts
var kafkaProducer sarama.SyncProducer
var assetCC *fabric.AssetCC

func main() {
	fmt.Println("Main Crosschain Service")
	var err error
	ethClient, err = ethclient.Dial(fmt.Sprintf("http://%s:8545", "localhost"))
	check(err)
	quorumClient, err = ethclient.Dial(fmt.Sprintf("http://%s:8546", "localhost"))
	check(err)

	assetCC = fabric.NewAssetCC()

	setupEthTransactor()

	zkNodes := []string{"localhost:2181"}
	kbConn, err := kb.NewConn(zkNodes)
	check(err)
	brokerList, _, err := kbConn.GetW()
	check(err)

	setupKafkaProducer(brokerList)
	go runKafkaConsumer(zkNodes)

	e := echo.New()
	e.Use(middleware.Recover())

	e.POST("/auction", handleCreateAuction)
	e.Logger.Fatal(e.Start(":9000"))
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

func runKafkaConsumer(zkNodes []string) {
	config := consumergroup.NewConfig()
	config.Offsets.Initial = sarama.OffsetNewest
	config.Offsets.ProcessingTimeout = 10 * time.Second

	consumer, err := consumergroup.JoinConsumerGroup(kafkaConsumerGroup,
		[]string{events.TopicOnEndAuction}, zkNodes, config)
	check(err)

	log.Printf("Subscribing %s", events.TopicOnEndAuction)

	for message := range consumer.Messages() {
		if message.Topic == events.TopicOnEndAuction {
			var event events.OnEndAuction
			err = json.Unmarshal(message.Value, &event)
			if err != nil {
				log.Printf("failed to parse event %+v", err)
				continue
			}
			log.Println("Received kafka event: OnEndAuction")
			go onEndAuction(event)
			consumer.CommitUpto(message)
		}
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func onEndAuction(event events.OnEndAuction) {
	args := fabric.EndAuctionArgs{
		AssetID: event.AssetID,
		AuctionResult: fabric.AuctionResult{
			Auction: fabric.Auction{
				ID: event.AuctionID,
			},
			HighestBidder: event.HighestBidder,
		},
	}
	_, err := assetCC.EndAuction(args)
	if err != nil {
		log.Printf("failed to end auction on fabric %+v", err)
		return
	}
	log.Printf("Ended auction on fabric")
}

func setupEthTransactor() {
	f, err := os.Open(privKeyFile)
	check(err)
	defer f.Close()
	ethTransactor, err = bind.NewTransactor(f, "password")
	check(err)
	ethTransactor.GasLimit = 1000000000
}
