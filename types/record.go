package types

import "time"

type Record struct {
	Id        string    `bun:"id,pk" json:"id"`
	Type      string    `bun:"type" json:"type"`
	User      string    `bun:"wallet_address" json:"wallet_address"`
	Amount    float64   `bun:"amount" json:"amount"`
	Currency  string    `bun:"currency" json:"currency"`
	CreatedAt time.Time `bun:"created_at" json:"created_at"`
}
