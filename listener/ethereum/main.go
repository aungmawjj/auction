package main

import (
	"auction/events"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wvanbergen/kafka/consumergroup"

	kb "github.com/philipjkim/kafka-brokers-go"
)

const kafkaConsumerGroup = "eth"

var ethClient *ethclient.Client
var kafkaProducer sarama.SyncProducer

func main() {
	fmt.Println("Ethereum listener")
	var err error
	ethClient, err = ethclient.Dial(fmt.Sprintf("http://%s:8545", "localhost"))
	check(err)

	zkNodes := []string{"localhost:2181"}
	kbConn, err := kb.NewConn(zkNodes)
	check(err)
	brokerList, _, err := kbConn.GetW()
	check(err)

	setupKafkaProducer(brokerList)
	runKafkaConsumer(zkNodes)
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
		[]string{events.TopicOnBindAuction}, zkNodes, config)
	check(err)

	log.Printf("Subscribing %s", events.TopicOnBindAuction)

	for message := range consumer.Messages() {
		if message.Topic == events.TopicOnBindAuction {
			var event events.OnBindAuction
			err = json.Unmarshal(message.Value, &event)
			if err != nil {
				log.Printf("failed to parse event %+v", err)
				continue
			}
			log.Printf("Received kafka event OnBindAuction")
			go onBindAuction(event)
			consumer.CommitUpto(message)
		}
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func onBindAuction(event events.OnBindAuction) {
	var err error
	result, err := listenAuctionEnd(event.AuctionID)
	if err != nil {
		log.Printf("failed to listen auction %+v", err)
	}
	time.Sleep(5 * time.Second)

	message := &sarama.ProducerMessage{Topic: events.TopicOnEndAuction, Partition: 0}
	value, _ := json.Marshal(events.OnEndAuction{
		AssetCC:       event.AssetCC,
		AssetID:       event.AssetID,
		AuctionID:     event.AuctionID,
		HighestBidder: result.HighestBidder,
	})
	message.Value = sarama.ByteEncoder(value)

	_, _, err = kafkaProducer.SendMessage(message)
	if err != nil {
		log.Printf("failed to send kafka message %+v", err)
		return
	}
	log.Print("Published event: OnEndAuction")
}
