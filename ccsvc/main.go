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
		[]string{events.TopicOnStartAuction, events.TopicOnEndAuction}, zkNodes, config)
	check(err)

	log.Printf("Subscribing %v", []string{events.TopicOnStartAuction, events.TopicOnEndAuction})

	for message := range consumer.Messages() {
		if message.Topic == events.TopicOnEndAuction {
			var payload events.AuctionEventPayload
			err = json.Unmarshal(message.Value, &payload)
			if err != nil {
				log.Printf("failed to parse event %+v", err)
				continue
			}
			log.Println("Received kafka event: ", events.TopicOnEndAuction)
			go onEndAuction(payload.AuctionID)
			consumer.CommitUpto(message)

		} else if message.Topic == events.TopicOnStartAuction {
			var payload events.AuctionEventPayload
			err = json.Unmarshal(message.Value, &payload)
			if err != nil {
				log.Printf("failed to parse event %+v", err)
				continue
			}
			log.Println("Received kafka event: ", events.TopicOnStartAuction)
			go onStartAuction(payload.AuctionID)
			consumer.CommitUpto(message)
		}
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func onStartAuction(auctionID int) {
	addrs, err := deployCrossChainAuctions()
	check(err)
	args := fabric.BindAuctionArgs{
		AuctionID:       auctionID,
		CrossAuctionIDs: addrs,
	}

	fmt.Println("Binding cross-chain auctions")
	_, err = assetCC.BindAuction(args)
	check(err)
}

func onEndAuction(auctionID int) {
	auction, err := assetCC.GetAuction(auctionID)
	check(err)

	highestBids := make([]int, 2)
	highestBidders := make([]string, 2)

	highestBids[0], highestBidders[0] = getAuctionInfo(auction.CrossAuctionIDs[0], ethClient)
	highestBids[1], highestBidders[1] = getAuctionInfo(auction.CrossAuctionIDs[1], quorumClient)

	args := fabric.EndAuctionArgs{
		AuctionID:      auctionID,
		HighestBids:    highestBids,
		HighestBidders: highestBidders,
	}

	_, err = assetCC.EndAuction(args)
	check(err)
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
