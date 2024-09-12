package main

import (
	"betcube_engine/config"
	"betcube_engine/services/nfticket_service/api"
	"betcube_engine/services/nfticket_service/db"

	datastore "betcube_engine/datastore"

	"github.com/gofiber/fiber/v2"
)

func main() {
	conf := config.NewConfig()
	nfticketStore := db.NewMixedNfticketDBStore(conf.PostgresDB, conf.MongoClient)

	app := fiber.New()
	nfticketHandler := api.NewNfticketHandler(&datastore.SupabaseStore{Nfticket: nfticketStore})

	// Set up routes
	apiv1 := app.Group("/v1")
	apiv1.Post("/nfticket", nfticketHandler.HandlePostNfticket)
	apiv1.Put("/nfticket/:nft_identifier", nfticketHandler.HandlePutNfticket)
	apiv1.Put("/user/nfticket/:nft_identifier", nfticketHandler.HandlePutNfticketByUser)
	apiv1.Get("/nftickets/:wallet_address", nfticketHandler.HandleGetNfticketsByWalletAddress)
	apiv1.Get("/nfticket/:nft_identifier", nfticketHandler.HandleGetNfticket)
}
