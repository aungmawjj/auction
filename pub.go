package main

import (
	"auction/contract"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
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

var (
	zkServers = flag.String("zk", os.Getenv("ZK_SERVERS"), "The comma-separated list of ZooKeeper servers. You can skip this flag by setting ZK_SERVERS environment variable")
	topic     = flag.String("topic", defaultKafkaTopic, "The topic to produce to")
	key       = flag.String("key", "", "The key of the message to produce. Can be empty.")
	silent    = flag.Bool("silent", false, "Turn off printing the message's topic, partition, and offset to stdout")

	logger = log.New(os.Stderr, "", log.LstdFlags)
)

func PublishAuctionResult(result *AuctionResult) error {
	flag.Parse()

	if *zkServers == "" {
		log.Fatalln("no -zk specified. Alternatively, set the ZK_SERVERS environment variable")
	}

	conn, err := kb.NewConn(strings.Split(*zkServers, ","))
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
			logger.Println("Failed to close Kafka producer cleanly:", err)
		}
	}()

	fmt.Println("Type a message and press Enter key to produce it. CTRL+C to exit.")

	message := &sarama.ProducerMessage{Topic: *topic, Partition: int32(0)}

	if *key != "" {
		message.Key = sarama.StringEncoder(*key)
	}

	value, _ := json.Marshal(result)
	message.Value = sarama.ByteEncoder(value)

	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		log.Fatalln(err)
	} else if !*silent {
		fmt.Printf("topic=%s\tpartition=%d\toffset=%d\n", *topic, partition, offset)
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
