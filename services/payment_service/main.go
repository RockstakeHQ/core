package main

import (
	"betcube_engine/config"
	"betcube_engine/services/payment_service/api"
	"betcube_engine/services/payment_service/db"

	datastore "betcube_engine/datastore"

	"github.com/gofiber/fiber/v2"
)

func main() {
	conf := config.NewConfig()
	paymentStore := db.NewStripeStore(conf.StripeKey)

	app := fiber.New()
	paymentHandler := api.NewPaymentHandler(&datastore.StripeStore{Payment: paymentStore})

	// Set up routes
	apiv1 := app.Group("/v1")
	apiv1.Post("/payment", paymentHandler.CreatePaymentIntent)
}
