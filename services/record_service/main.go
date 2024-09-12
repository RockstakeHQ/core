package main

import (
	"betcube_engine/config"
	"betcube_engine/services/record_service/api"
	"betcube_engine/services/record_service/db"

	datastore "betcube_engine/datastore"

	"github.com/gofiber/fiber/v2"
)

func main() {
	conf := config.NewConfig()
	recordStore := db.NewSupabaseRecordStore(conf.PostgresDB)
	app := fiber.New()
	recordHandler := api.NewRecordHandler(&datastore.SupabaseStore{Record: recordStore})

	api := app.Group("/v1")
	api.Post("/record", recordHandler.HandlePostRecord)
	api.Get("/records/:wallet_address", recordHandler.HandleGetRecordsByUser)

	app.Listen(":4000")
}
