package events

import (
	"checkout/src/database"
	"checkout/src/models"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func Listen(message *kafka.Message) error {
	database.Connect()

	key := string(message.Key)

	switch key {
	case "link_created":
		return LinkCreated(message.Value)
	case "product_created":
		return ProductCreated(message.Value)
	case "product_updated":
		return ProductUpdated(message.Value)
	case "product_deleted":
		return ProductDeleted(message.Value)
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

func ProductCreated(value []byte) error {
	var product models.Product

	if err := json.Unmarshal(value, &product); err != nil {
		return err
	}

	if err := database.DB.Create(&product).Error; err != nil {
		return err
	}

	return nil
}

func ProductUpdated(value []byte) error {
	var product models.Product

	if err := json.Unmarshal(value, &product); err != nil {
		return err
	}

	if err := database.DB.Model(&product).Updates(&product).Error; err != nil {
		return err
	}

	return nil
}

func ProductDeleted(value []byte) error {
	var id uint

	if err := json.Unmarshal(value, &id); err != nil {
		return err
	}

	if err := database.DB.Delete(&models.Product{}, "id = ?", id).Error; err != nil {
		return err
	}

	return nil
}
