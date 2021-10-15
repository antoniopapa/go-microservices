package events

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"context"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func Listen(message *kafka.Message) error {
	key := string(message.Key)

	switch key {
	case "product_created":
		return ProductCreated(message.Value)
	case "product_updated":
		return ProductUpdated(message.Value)
	case "product_deleted":
		return ProductDeleted(message.Value)
	case "order_created":
		return OrderCreated(message.Value)
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

	go database.ClearCache("products_frontend", "products_backend")

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

	go database.ClearCache("products_frontend", "products_backend")

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

	go database.ClearCache("products_frontend", "products_backend")

	return nil
}

type Order struct {
	Id                uint    `json:"id"`
	UserId            uint    `json:"user_id"`
	Code              string  `json:"code"`
	AmbassadorRevenue float64 `json:"ambassador_revenue"`
	AmbassadorName    string  `json:"ambassador_name"`
}

func OrderCreated(value []byte) error {
	var order Order

	if err := json.Unmarshal(value, &order); err != nil {
		return err
	}

	newOrder := models.Order{
		Id:     order.Id,
		UserId: order.UserId,
		Code:   order.Code,
		Total:  order.AmbassadorRevenue,
	}

	if err := database.DB.Create(&newOrder).Error; err != nil {
		return err
	}

	database.Cache.ZIncrBy(context.Background(), "rankings", order.AmbassadorRevenue, order.AmbassadorName)

	return nil
}
