package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
	"users/src/database"
	"users/src/models"
)

const SecretKey = "secret"

type ClaimsWithScope struct {
	jwt.StandardClaims
	Scope string
}

func IsAuthenticated(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &ClaimsWithScope{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil || !token.Valid {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	payload := token.Claims.(*ClaimsWithScope)

	id, _ := strconv.Atoi(payload.Subject)

	var user models.User

	database.DB.Where("id = ?", id).First(&user)

	if user.Id == 0 {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	var userToken models.UserToken

	database.DB.Where("user_id = ? and token = ? and expired_at >= now()", id, token.Raw).Last(&userToken)

	if userToken.Id == 0 {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	c.Context().SetUserValue("scope", payload.Scope)
	c.Context().SetUserValue("user", user)

	return c.Next()
}

func GenerateJWT(id uint, scope string) (string, error) {
	payload := ClaimsWithScope{}
	payload.Subject = strconv.Itoa(int(id))
	payload.ExpiresAt = time.Now().Add(time.Hour * 24).Unix()
	payload.Scope = scope

	return jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte(SecretKey))
}
