package main

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"net/http"
	"time"
)

type TestReq struct {
	Body  []interface{} `json:"body"`
	Delay int64         `json:"delay"`
}

func main() {

	http.HandleFunc("/test", func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("test tesst")
		decoder := json.NewDecoder(req.Body)
		var t TestReq
		err := decoder.Decode(&t)
		if err != nil {
			panic(err)
		}

		// to produce messages
		topic := "oms.order"

		p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
		if err != nil {
			panic(err)
		}

		// Delivery report handler for produced messages
		go func() {
			for e := range p.Events() {
				switch ev := e.(type) {
				case *kafka.Message:
					if ev.TopicPartition.Error != nil {
						fmt.Printf("Delivery failed: %v\n", ev.Value)
					} else {
						fmt.Printf("Delivered message to %v\n", ev.Value)
					}
				}
			}
		}()

		// Produce messages to topic (asynchronously)
		for _, message := range t.Body {
			data, err := json.Marshal(message)
			if err != nil {
				fmt.Println("error:", err)
			}
			p.Produce(&kafka.Message{
				TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
				Value:          data,
			}, nil)

			time.Sleep(time.Duration(t.Delay) * time.Millisecond)
		}

		// Wait for message deliveries
		p.Flush(15 * 1000)
	})
	log.Printf("Listening on localhost:8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
