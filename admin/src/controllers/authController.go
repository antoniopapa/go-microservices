package controllers

import (
	"admin/src/models"
	"admin/src/services"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"time"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	data["is_ambassador"] = "false"

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

	data["scope"] = "admin"

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
	return c.JSON(c.Context().UserValue("user"))
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

	c.Context().SetUserValue("user", "")

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
