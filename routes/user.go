package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sajjadanwar0/clubhouse-clone/config"
	"github.com/sajjadanwar0/clubhouse-clone/handlers"
	"github.com/sajjadanwar0/clubhouse-clone/middleware"
)

func SetupUserRoutes(grp fiber.Router, handlers *handlers.Handler) {
	conf := config.New()
	useRoute := grp.Group("/user")
	useRoute.Post("/register", handlers.UserRegister)
	useRoute.Post("/login", handlers.UserLogin)
	useRoute.Use(middleware.IsAuthenticated(conf))
	useRoute.Post("/me", handlers.MeQuery)

}
