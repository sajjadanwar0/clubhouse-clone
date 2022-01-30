package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sajjadanwar0/clubhouse-clone/handlers"
)

func SetUpApiV1(app *fiber.App, handlers *handlers.Handler) {

	v1 := app.Group("/api/v1")
	SetupUserRoutes(v1, handlers)

}
