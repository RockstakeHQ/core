package api

import (
	datastore "betcube_engine/datastore"

	"github.com/gofiber/fiber/v2"
)

type PaymentHandler struct {
	store *datastore.StripeStore
}

func NewPaymentHandler(s *datastore.StripeStore) *PaymentHandler {
	return &PaymentHandler{
		store: s,
	}
}

func (h *PaymentHandler) CreatePaymentIntent(c *fiber.Ctx) error {
	var body struct {
		Amount   int64  `json:"amount"`
		Currency string `json:"currency"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}

	paymentIntent, err := h.store.Payment.CreateStripePaymentIntent(body.Amount, body.Currency)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if paymentIntent.Status != "succeeded" {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message":       "Confirm payment please",
			"client_secret": paymentIntent.ClientSecret,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Payment Completed Successfully",
	})
}
