package main

import (
	"betcube-engine/config"
	"betcube-engine/services/record_service/api"
	"betcube-engine/services/record_service/db"

	datastore "betcube-engine/datastore"

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
