package main

import (
	"betcube-engine/config"
	"betcube-engine/services/user_service/api"
	"betcube-engine/services/user_service/db"

	datastore "betcube-engine/datastore"

	"github.com/gofiber/fiber/v2"
)

func main() {
	conf := config.NewConfig()
	userStore := db.NewSupabaseUserStore(conf.PostgresDB)
	app := fiber.New()
	userHandler := api.NewUserHandler(&datastore.SupabaseStore{User: userStore})

	api := app.Group("/v1")
	api.Post("/user", userHandler.HandlePostUser)
	api.Get("/user/wallet/:wallet_address", userHandler.HandleGetUserByWalletAddress)
	api.Get("/user/id/:user_id", userHandler.HandleGetUserByUserId)

	app.Listen(":4000")
}
