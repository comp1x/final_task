package main

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gitlab.com/mediasoft-internship/final-task/contracts/pkg/contracts/customer"
	"google.golang.org/protobuf/proto"
	"log"
	"strconv"
	"time"
)

type OrderPlacer struct {
	producer   *kafka.Producer
	topic      string
	deliverych chan kafka.Event
}

func NewOrderPlacer(p *kafka.Producer, topic string) *OrderPlacer {
	return &OrderPlacer{
		producer:   p,
		topic:      topic,
		deliverych: make(chan kafka.Event, 10000),
	}
}

func (op *OrderPlacer) placeOrder(request customer.CreateOrderRequest) error {
	var (
		payload, err = proto.Marshal(&request)
	)
	if err != nil {
		log.Fatal(err)
	}
	err = op.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &op.topic,
			Partition: kafka.PartitionAny,
		},
		Value: payload,
	},
		op.deliverych,
	)
	if err != nil {
		log.Fatal(err)
	}
	<-op.deliverych
	return nil
}

func main() {

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"client.id":         "orders",
		"acks":              "all",
	})
	if err != nil {
		log.Fatal("failed when create kafka producer: ", err)
	}

	jsonPayload := `{
  		"meats": [
    	{
      		"product_uuid": "d54b2d5f-6403-4a7a-ac73-e66bb392caf7",
      		"count": 322
		}
  		],
  		"user_uuid": "ce1b22aa-6720-4cf6-b3e8-bec68f07c1de"
	}`

	var request customer.CreateOrderRequest

	if err := json.Unmarshal([]byte(jsonPayload), &request); err != nil {
		log.Fatal("failed with unmarshalling request: ", err)
	}

	op := NewOrderPlacer(p, "orders")
	for i := 0; i < 1000; i++ {
		if err := op.placeOrder(request); err != nil {
			log.Fatal("failed when place order: ", err)
		}
		fmt.Println("added in kafka req №" + strconv.Itoa(i))
		time.Sleep(time.Millisecond * 100)
	}
}
