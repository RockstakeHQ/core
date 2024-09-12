package main

import (
	"betcube-engine/config"
	"betcube-engine/services/sportsbook_service/api"
	"betcube-engine/services/sportsbook_service/db"
	"os"

	datastore "betcube-engine/datastore"

	"github.com/gofiber/fiber/v2"
)

func main() {
	conf := config.NewConfig()
	sportsbookStore := db.NewMongoSportsbookStore(conf.MongoClient)

	app := fiber.New()
	sportsbookHandler := api.NewSportsbookHandler(&datastore.MongoStore{Sportsbook: sportsbookStore})

	apiv1Football := app.Group("/v1/football")
	apiv1Football.Post("/sportsbook", sportsbookHandler.HandlePostSportsbook)
	apiv1Football.Get("/sportsbook/:date", sportsbookHandler.HandleGetSportsbookByDate)
	apiv1Football.Get("/sportsbook/fixture/:id", sportsbookHandler.HandlerGetSportsbookByFixtureID)

	listenAddr := os.Getenv("HTTP_LISTEN_ADDRESS")
	app.Listen(listenAddr)
}
