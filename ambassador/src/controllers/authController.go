package controllers

import (
	"ambassador/src/database"
	"ambassador/src/models"
	"ambassador/src/services"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"time"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	data["is_ambassador"] = "true"

	response, err := services.UserService.Post("register", "", data)

	if err != nil {
		return err
	}

	var user models.User

	json.NewDecoder(response.Body).Decode(&user)

	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	data["scope"] = "ambassador"

	response, err := services.UserService.Post("login", "", data)

	if err != nil {
		return err
	}

	var res map[string]string

	json.NewDecoder(response.Body).Decode(&res)

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    res["jwt"],
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func User(c *fiber.Ctx) error {
	ambassador := models.Ambassador(c.Context().UserValue("user").(models.User))

	ambassador.CalculateRevenue(database.DB)

	return c.JSON(ambassador)
}

func Logout(c *fiber.Ctx) error {
	services.UserService.Post("logout", c.Cookies("jwt", ""), nil)

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func UpdateInfo(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	response, err := services.UserService.Put("users/info", c.Cookies("jwt", ""), data)

	if err != nil {
		return err
	}

	var user models.User

	json.NewDecoder(response.Body).Decode(&user)

	return c.JSON(user)
}

func UpdatePassword(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	response, err := services.UserService.Put("users/password", c.Cookies("jwt", ""), data)

	if err != nil {
		return err
	}

	var user models.User

	json.NewDecoder(response.Body).Decode(&user)

	return c.JSON(user)
}
