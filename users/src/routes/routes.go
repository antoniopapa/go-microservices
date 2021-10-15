package routes

import (
	"github.com/gofiber/fiber/v2"
	"users/src/controllers"
	"users/src/middlewares"
)

func Setup(app *fiber.App) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON("ok")
	})

	api := app.Group("api")

	api.Post("register", controllers.Register)
	api.Post("login", controllers.Login)
	api.Get("users", controllers.Users)
	api.Get("users/:id", controllers.GetUser)

	authenticated := api.Use(middlewares.IsAuthenticated)
	api.Get("user/:scope", controllers.User)
	authenticated.Post("logout", controllers.Logout)
	authenticated.Put("users/info", controllers.UpdateInfo)
	authenticated.Put("users/password", controllers.UpdatePassword)
}
