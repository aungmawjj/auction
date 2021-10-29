package main

import (
	"auction/contract"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	kb "github.com/philipjkim/kafka-brokers-go"
)

type AuctionResult struct {
	Address       []byte
	Ended         bool
	HighestBid    int64
	HighestBidder []byte
}

func fetchAuctionResult(address []byte) error {
	eth, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		return err
	}
	cc, err := contract.NewAuction(common.BytesToAddress(address), eth)
	if err != nil {
		return err
	}
	var result AuctionResult
	result.Address = address
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

	return publishAuctionResult(&result)
}

const (
	defaultKafkaTopic = "test_topic"
)

func PublishAuctionResult(result *AuctionResult) error {
	flag.Parse()

	conn, err := kb.NewConn(strings.Split("localhost:2181", ","))
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	brokerList, _, err := conn.GetW()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("brokerList: %q\n", brokerList)

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewManualPartitioner

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln(err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Println("Failed to close Kafka producer cleanly:", err)
		}
	}()

	fmt.Println("Type a message and press Enter key to produce it. CTRL+C to exit.")

	message := &sarama.ProducerMessage{Topic: defaultKafkaTopic, Partition: int32(0)}

	value, _ := json.Marshal(result)
	message.Value = sarama.ByteEncoder(value)

	_, _, err = producer.SendMessage(message)
	if err != nil {
		log.Fatalln(err)
	}
	return nil
}

func publishAuctionResult(result *AuctionResult) error {
	buf := bytes.NewBuffer(nil)
	json.NewEncoder(buf).Encode(result)

	resp, err := http.Post("http://localhost:9000/on_receive", "application/json", buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to submit auction result, status: %d", resp.StatusCode)
	}

	return nil
}
