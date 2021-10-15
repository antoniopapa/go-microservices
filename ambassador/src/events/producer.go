package events

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
)

var Producer *kafka.Producer

func SetupProducer() {
	var err error

	Producer, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("BOOTSTRAP_SERVERS"),
		"security.protocol": os.Getenv("SECURITY_PROTOCOL"),
		"sasl.username":     os.Getenv("SASL_USERNAME"),
		"sasl.password":     os.Getenv("SASL_PASSWORD"),
		"sasl.mechanism":    os.Getenv("SASL_MECHANISM"),
	})

	if err != nil {
		panic(err)
	}

	//defer Producer.Close()
}

func Produce(topic string, key string, message interface{}) {
	value, _ := json.Marshal(message)

	Producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            []byte(key),
		Value:          value,
	}, nil)

	Producer.Flush(15000)
}
