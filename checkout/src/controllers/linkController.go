package controllers

import (
	"checkout/src/database"
	"checkout/src/models"
	"checkout/src/services"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func GetLink(c *fiber.Ctx) error {
	code := c.Params("code")

	var link models.Link

	database.DB.Preload("Products").Where("code = ?", code).First(&link)

	if link.Id == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid link!",
		})
	}

	response, err := services.UserService.Get(fmt.Sprintf("users/%d", link.UserId), "")

	if err != nil {
		return err
	}

	var user models.User

	json.NewDecoder(response.Body).Decode(&user)

	link.User = user

	return c.JSON(link)
}
