package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"github.com/sajjadanwar0/clubhouse-clone/config"
	"github.com/sajjadanwar0/clubhouse-clone/ent"
	"github.com/sajjadanwar0/clubhouse-clone/ent/migrate"
	handlers2 "github.com/sajjadanwar0/clubhouse-clone/handlers"
	"github.com/sajjadanwar0/clubhouse-clone/middleware"
	"github.com/sajjadanwar0/clubhouse-clone/routes"
	"github.com/sajjadanwar0/clubhouse-clone/utils"
	"log"
)

func main() {

	conf := config.New()
	client, err := ent.Open("postgres", fmt.Sprintf("host=%s port=%s dbname=%s password=%s sslmode=disable", conf.Database.Host, conf.Database.Port, conf.Database.Name, conf.Database.Password))

	if err != nil {
		utils.Fataf("Database connection failed :", err)
	}
	defer client.Close()

	ctx := context.Background()
	err = client.Schema.Create(ctx, migrate.WithDropIndex(true), migrate.WithDropColumn(true))
	if err != nil {
		utils.Fataf("Migration fail :", err)
	}
	app := fiber.New()
	middleware.SetMiddleware(app)
	handler := handlers2.NewHandlers(client, conf)
	routes.SetUpApiV1(app, handler)

	port := "4000"

	addr := flag.String("addr", port, "http service address")
	flag.Parse()
	log.Fatalln(app.Listen(":" + *addr))
}
