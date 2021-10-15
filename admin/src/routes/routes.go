package routes

import (
	"admin/src/controllers"
	"admin/src/middlewares"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON("ok")
	})
	
	api := app.Group("api/admin")
	api.Post("register", controllers.Register)
	api.Post("login", controllers.Login)

	adminAuthenticated := api.Use(middlewares.IsAuthenticated)
	adminAuthenticated.Get("user", controllers.User)
	adminAuthenticated.Post("logout", controllers.Logout)
	adminAuthenticated.Put("users/info", controllers.UpdateInfo)
	adminAuthenticated.Put("users/password", controllers.UpdatePassword)
	adminAuthenticated.Get("ambassadors", controllers.Ambassadors)
	adminAuthenticated.Get("products", controllers.Products)
	adminAuthenticated.Post("products", controllers.CreateProducts)
	adminAuthenticated.Get("products/:id", controllers.GetProduct)
	adminAuthenticated.Put("products/:id", controllers.UpdateProduct)
	adminAuthenticated.Delete("products/:id", controllers.DeleteProduct)
	adminAuthenticated.Get("users/:id/links", controllers.Link)
	adminAuthenticated.Get("orders", controllers.Orders)
}
