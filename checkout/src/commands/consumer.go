package main

import (
	"checkout/src/database"
	"checkout/src/events"
	"checkout/src/models"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
)

func main() {
	database.Connect()

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("BOOTSTRAP_SERVERS"),
		"security.protocol": os.Getenv("SECURITY_PROTOCOL"),
		"sasl.username":     os.Getenv("SASL_USERNAME"),
		"sasl.password":     os.Getenv("SASL_PASSWORD"),
		"sasl.mechanism":    os.Getenv("SASL_MECHANISM"),
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		panic(err)
	}

	consumer.SubscribeTopics([]string{os.Getenv("KAFKA_TOPIC")}, nil)

	fmt.Println("Started Consuming")

	for {
		msg, err := consumer.ReadMessage(-1)

		if err != nil {
			// The client will automatically try to recover from all errors.
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)

			database.DB.Create(&models.KafkaError{
				Key:   msg.Key,
				Value: msg.Value,
				Error: err,
			})

			return
		}

		fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))

		if err := events.Listen(msg); err != nil {
			database.DB.Create(&models.KafkaError{
				Key:   msg.Key,
				Value: msg.Value,
				Error: err,
			})
		}
	}

	consumer.Close()
}
