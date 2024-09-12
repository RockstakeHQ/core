package main

import (
	"betcube_engine/config"
	"betcube_engine/services/sportsbook_service/api"
	"betcube_engine/services/sportsbook_service/db"
	"os"

	datastore "betcube_engine/datastore"

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
