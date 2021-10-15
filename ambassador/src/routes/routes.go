package routes

import (
	"ambassador/src/controllers"
	"ambassador/src/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON("ok")
	})

	api := app.Group("api/ambassador")

	api.Post("register", controllers.Register)
	api.Post("login", controllers.Login)
	api.Get("products/frontend", controllers.ProductsFrontend)
	api.Get("products/backend", controllers.ProductsBackend)

	authenticated := api.Use(middlewares.IsAuthenticated)
	authenticated.Get("user", controllers.User)
	authenticated.Post("logout", controllers.Logout)
	authenticated.Put("users/info", controllers.UpdateInfo)
	authenticated.Put("users/password", controllers.UpdatePassword)
	authenticated.Post("links", controllers.CreateLink)
	authenticated.Get("stats", controllers.Stats)
	authenticated.Get("rankings", controllers.Rankings)
}
