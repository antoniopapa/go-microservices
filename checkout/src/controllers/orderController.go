package controllers

import (
	"checkout/src/database"
	"checkout/src/events"
	"checkout/src/models"
	"checkout/src/services"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"os"
)

type CreateOrderRequest struct {
	Code      string
	FirstName string
	LastName  string
	Email     string
	Address   string
	Country   string
	City      string
	Zip       string
	Products  []map[string]int
}

func CreateOrder(c *fiber.Ctx) error {
	var request CreateOrderRequest

	if err := c.BodyParser(&request); err != nil {
		return err
	}

	link := models.Link{
		Code: request.Code,
	}

	database.DB.First(&link)

	response, err := services.UserService.Get(fmt.Sprintf("users/%d", link.UserId), "")

	if err != nil {
		return err
	}

	var user models.User

	json.NewDecoder(response.Body).Decode(&user)

	if link.Id == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid link!",
		})
	}

	order := models.Order{
		Code:            link.Code,
		UserId:          link.UserId,
		AmbassadorEmail: user.Email,
		FirstName:       request.FirstName,
		LastName:        request.LastName,
		Email:           request.Email,
		Address:         request.Address,
		Country:         request.Country,
		City:            request.City,
		Zip:             request.Zip,
	}

	tx := database.DB.Begin()

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var lineItems []*stripe.CheckoutSessionLineItemParams

	for _, requestProduct := range request.Products {
		product := models.Product{}
		product.Id = uint(requestProduct["product_id"])
		database.DB.First(&product)

		total := product.Price * float64(requestProduct["quantity"])

		item := models.OrderItem{
			OrderId:           order.Id,
			ProductTitle:      product.Title,
			Price:             product.Price,
			Quantity:          uint(requestProduct["quantity"]),
			AmbassadorRevenue: 0.1 * total,
			AdminRevenue:      0.9 * total,
		}

		if err := tx.Create(&item).Error; err != nil {
			tx.Rollback()
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			Name:        stripe.String(product.Title),
			Description: stripe.String(product.Description),
			Images:      []*string{stripe.String(product.Image)},
			Amount:      stripe.Int64(100 * int64(product.Price)),
			Currency:    stripe.String("usd"),
			Quantity:    stripe.Int64(int64(requestProduct["quantity"])),
		})
	}

	stripe.Key = os.Getenv("STRIPE_KEY")

	params := stripe.CheckoutSessionParams{
		SuccessURL:         stripe.String(fmt.Sprintf("%s/success?source={CHECKOUT_SESSION_ID}", os.Getenv("CHECKOUT_URL"))),
		CancelURL:          stripe.String(fmt.Sprintf("%s/error", os.Getenv("CHECKOUT_URL"))),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems:          lineItems,
	}

	source, err := session.New(&params)

	if err != nil {
		tx.Rollback()
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	order.TransactionId = source.ID

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	tx.Commit()

	return c.JSON(source)
}

func CompleteOrder(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	order := models.Order{}

	database.DB.Preload("OrderItems").First(&order, models.Order{
		TransactionId: data["source"],
	})

	if order.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "Order not found",
		})
	}

	order.Complete = true
	database.DB.Save(&order)

	ambassadorRevenue := 0.0
	adminRevenue := 0.0

	for _, item := range order.OrderItems {
		ambassadorRevenue += item.AmbassadorRevenue
		adminRevenue += item.AdminRevenue
	}

	response, err := services.UserService.Get(fmt.Sprintf("users/%d", order.UserId), "")

	if err != nil {
		return err
	}

	var user models.User

	json.NewDecoder(response.Body).Decode(&user)

	order.AmbassadorName = user.Name()
	order.AmbassadorRevenue = ambassadorRevenue
	order.AdminRevenue = adminRevenue

	go events.Produce("admin_topic", "order_created", order)
	go events.Produce("ambassador_topic", "order_created", order)
	go events.Produce("email_topic", "order_created", order)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}
