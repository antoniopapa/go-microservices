package routes

import (
	"checkout/src/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON("ok")
	})

	api := app.Group("api/checkout")
	api.Get("links/:code", controllers.GetLink)
	api.Post("orders", controllers.CreateOrder)
	api.Post("orders/confirm", controllers.CompleteOrder)
}
