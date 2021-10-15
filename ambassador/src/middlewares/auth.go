package middlewares

import (
	"ambassador/src/models"
	"ambassador/src/services"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
)

func IsAuthenticated(c *fiber.Ctx) error {
	response, err := services.UserService.Get("user/ambassador", c.Cookies("jwt", ""))

	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	var user models.User

	json.NewDecoder(response.Body).Decode(&user)

	if user.Id == 0 {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	c.Context().SetUserValue("user", user)

	return c.Next()
}
