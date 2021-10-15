package controllers

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
	"users/src/database"
	"users/src/middlewares"
	"users/src/models"
)

func Register(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "passwords do not match",
		})
	}

	isAmbassador, _ := strconv.ParseBool(data["is_ambassador"])

	user := models.User{
		FirstName:    data["first_name"],
		LastName:     data["last_name"],
		Email:        data["email"],
		IsAmbassador: isAmbassador,
	}
	user.SetPassword(data["password"])

	database.DB.Create(&user)

	return c.JSON(user)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.Id == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid Credentials",
		})
	}

	if err := user.ComparePassword(data["password"]); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid Credentials",
		})
	}

	if data["scope"] != "ambassador" && user.IsAmbassador {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	token, err := middlewares.GenerateJWT(user.Id, data["scope"])

	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid Credentials",
		})
	}

	userToken := models.UserToken{
		UserId:    user.Id,
		Token:     token,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(time.Hour * 24),
	}

	database.DB.Create(&userToken)

	return c.JSON(fiber.Map{
		"jwt": token,
	})
}

func User(c *fiber.Ctx) error {
	scopeAdminParamAmbassador := c.Context().UserValue("scope") == "admin" && c.Params("scope") == "ambassador"
	scopeAmbassadorParamAdmin := c.Context().UserValue("scope") == "ambassador" && c.Params("scope") == "admin"

	if scopeAdminParamAmbassador || scopeAmbassadorParamAdmin {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthorized",
		})
	}

	return c.JSON(c.Context().UserValue("user"))
}

func Logout(c *fiber.Ctx) error {
	user := c.Context().UserValue("user").(models.User)

	database.DB.Delete(models.UserToken{}, "user_id = ?", user.Id)

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

	user := c.Context().UserValue("user").(models.User)

	user.FirstName = data["first_name"]
	user.LastName = data["last_name"]
	user.Email = data["email"]

	database.DB.Model(&user).Updates(&user)

	return c.JSON(user)
}

func UpdatePassword(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "passwords do not match",
		})
	}

	user := c.Context().UserValue("user").(models.User)
	user.SetPassword(data["password"])

	database.DB.Model(&user).Updates(&user)

	return c.JSON(user)
}

func Users(c *fiber.Ctx) error {
	var users []models.User

	database.DB.Find(&users)

	return c.JSON(users)
}

func GetUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var user models.User
	user.Id = uint(id)

	database.DB.Find(&user)

	return c.JSON(user)
}
