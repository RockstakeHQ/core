package types

import "time"

type Market struct {
	ID         int       `json:"id"`
	FixtureID  int       `json:"fixture_id"`
	MarketType string    `json:"market_type"`
	Liquidity  float64   `json:"liquidity"`
	OrderBook  OrderBook `json:"order_book"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type OrderBook struct {
	Backs []Order `json:"backs"`
	Lays  []Order `json:"lays"`
}

type Order struct {
	Price    float64   `json:"price"`
	Amount   float64   `json:"amount"`
	DateTime time.Time `json:"datetime"`
}

type Bet struct {
	ID        int       `json:"id"`
	WalletID  string    `json:"wallet_id"`
	NFTID     int       `json:"nft_id"`
	MarketID  int       `json:"market_id"` // Referință la market-ul specific
	Side      string    `json:"side"`      // "back" sau "lay"
	Stake     float64   `json:"stake"`     // Miza
	Odds      float64   `json:"odds"`      // Cota
	Status    string    `json:"status"`    // ex: "pending", "matched", "settled"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
