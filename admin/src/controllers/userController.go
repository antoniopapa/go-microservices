package controllers

import (
	"admin/src/models"
	"admin/src/services"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
)

func Ambassadors(c *fiber.Ctx) error {
	response, err := services.UserService.Get("users", c.Cookies("jwt", ""))

	if err != nil {
		return err
	}

	var users []models.User

	ambassadors := []models.User{}

	json.NewDecoder(response.Body).Decode(&users)

	for _, user := range users {
		if user.IsAmbassador {
			ambassadors = append(ambassadors, user)
		}
	}

	return c.JSON(ambassadors)
}
