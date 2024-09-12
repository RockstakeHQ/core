package db

import (
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
)

type PaymentStore interface {
	CreateStripePaymentIntent(amount int64, currency string) (*stripe.PaymentIntent, error)
}

type Store struct {
	key string
}

func NewStripeStore(key string) PaymentStore {
	return &Store{
		key: key,
	}
}

func (s *Store) CreateStripePaymentIntent(amount int64, currency string) (*stripe.PaymentIntent, error) {
	stripe.Key = s.key

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(currency),
	}

	return paymentintent.New(params)
}
