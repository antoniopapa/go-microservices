package events

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var Producer *kafka.Producer

func SetupProducer() {
	var err error

	Producer, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "pkc-4ygn6.europe-west3.gcp.confluent.cloud:9092",
		"security.protocol": "SASL_SSL",
		"sasl.username":     "6567MKPS6YFGEMJH",
		"sasl.password":     "zaf8YnUY/BzmTpNk6TFC+wJQQ0YoQdvmwfRQTRtRo19YZChYEev2Wkx+rG3Two2h",
		"sasl.mechanism":    "PLAIN",
	})

	if err != nil {
		panic(err)
	}

	//defer Producer.Close()
}

func Produce(topic string, key string, message interface{}) {
	value, _ := json.Marshal(message)

	err := Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          value,
	}, nil)

	if err != nil {
		panic(err)
	}

	Producer.Flush(15000)
}
