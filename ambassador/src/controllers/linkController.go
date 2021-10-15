package controllers

import (
	"ambassador/src/database"
	"ambassador/src/events"
	"ambassador/src/models"
	"github.com/bxcodec/faker/v3"
	"github.com/gofiber/fiber/v2"
)

type CreateLinkRequest struct {
	Products []int
}

func CreateLink(c *fiber.Ctx) error {
	var request CreateLinkRequest

	if err := c.BodyParser(&request); err != nil {
		return err
	}

	user := c.Context().UserValue("user").(models.User)

	link := models.Link{
		UserId: user.Id,
		Code:   faker.Username(),
	}

	for _, productId := range request.Products {
		product := models.Product{}
		product.Id = uint(productId)
		link.Products = append(link.Products, product)
	}

	database.DB.Create(&link)

	go events.Produce("admin_topic", "link_created", link)
	go events.Produce("checkout_topic", "link_created", link)

	return c.JSON(link)
}

func Stats(c *fiber.Ctx) error {
	user := c.Context().UserValue("user").(models.User)

	var links []models.Link

	database.DB.Find(&links, models.Link{
		UserId: user.Id,
	})

	var result []interface{}

	var orders []models.Order

	for _, link := range links {
		database.DB.Find(&orders, &models.Order{
			Code: link.Code,
		})

		revenue := 0.0

		for _, order := range orders {
			revenue += order.Total
		}

		result = append(result, fiber.Map{
			"code":    link.Code,
			"count":   len(orders),
			"revenue": revenue,
		})
	}

	return c.JSON(result)
}
