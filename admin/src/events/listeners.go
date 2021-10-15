package events

import (
	"admin/src/database"
	"admin/src/models"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func Listen(message *kafka.Message) error {
	database.Connect()

	key := string(message.Key)

	switch key {
	case "link_created":
		return LinkCreated(message.Value)
	case "order_created":
		return OrderCreated(message.Value)
	}

	return nil
}

func LinkCreated(value []byte) error {
	var link models.Link

	if err := json.Unmarshal(value, &link); err != nil {
		return err
	}

	if err := database.DB.Create(&link).Error; err != nil {
		return err
	}

	return nil
}

func OrderCreated(value []byte) error {
	var order models.Order

	if err := json.Unmarshal(value, &order); err != nil {
		return err
	}

	if err := database.DB.Create(&order).Error; err != nil {
		return err
	}

	return nil
}
