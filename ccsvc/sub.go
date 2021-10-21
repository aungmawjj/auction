package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/labstack/echo/v4"
	kb "github.com/philipjkim/kafka-brokers-go"
	"github.com/wvanbergen/kafka/consumergroup"
)

const (
	defaultKafkaTopic    = "test_topic"
	defaultConsumerGroup = "defaultConsumerGroup"
)

func Subscribe() {
	flag.Parse()

	zkNodes := strings.Split("localhost", ",")
	conn, err := kb.NewConn(zkNodes)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	brokerList, _, err := conn.GetW()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("brokerList: %q\n", brokerList)

	config := consumergroup.NewConfig()
	config.Offsets.Initial = sarama.OffsetNewest
	config.Offsets.ProcessingTimeout = 10 * time.Second

	consumer, consumerErr := consumergroup.JoinConsumerGroup(defaultConsumerGroup,
		[]string{defaultKafkaTopic}, zkNodes, config)
	if consumerErr != nil {
		log.Fatalln(consumerErr)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		if err := consumer.Close(); err != nil {
			sarama.Logger.Println("Error closing the consumer", err)
		}
	}()

	go func() {
		for err := range consumer.Errors() {
			log.Println(err)
		}
	}()

	eventCount := 0
	offsets := make(map[string]map[int32]int64)

	for message := range consumer.Messages() {
		if offsets[message.Topic] == nil {
			offsets[message.Topic] = make(map[int32]int64)
		}

		eventCount++
		if offsets[message.Topic][message.Partition] != 0 &&
			offsets[message.Topic][message.Partition] != message.Offset-1 {
			log.Printf(
				"Unexpected offset on %s:%d. Expected %d, found %d, diff %d.\n",
				message.Topic,
				message.Partition,
				offsets[message.Topic][message.Partition]+1,
				message.Offset,
				message.Offset-offsets[message.Topic][message.Partition]+1)
		}

		log.Printf("partition: %d, offset: %d, key: %s, value: %s",
			message.Partition, message.Offset, message.Key, message.Value)

		var result AuctionResult
		err := json.Unmarshal(message.Value, &result)
		if err == nil {
			onReceiveAuctionResult(&result)
		}

		time.Sleep(10 * time.Millisecond)

		offsets[message.Topic][message.Partition] = message.Offset
		consumer.CommitUpto(message)
	}

	log.Printf("Processed %d events", eventCount)
	log.Printf("%+v", offsets)
}

func handleOnReceiveAuctionResult(c echo.Context) error {

	var result AuctionResult
	err := c.Bind(&result)
	if err != nil {
		return err
	}
	onReceiveAuctionResult(&result)

	return c.JSON(http.StatusOK, result)
}
