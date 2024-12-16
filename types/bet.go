package types

import "time"

type Bet struct {
	ID            string    `json:"id"`
	WalletAddress string    `json:"wallet_address"`
	FixtureID     int       `json:"fixture_id"`
	MarketID      int       `json:"market_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Status        string    `json:"status"`
	BetInfo       BetInfo   `json:"bet_info"`
}

type BetInfo struct {
	Selection    string  `json:"selection"`
	Type         string  `json:"type"`
	Odd          float64 `json:"odd"`
	Stake        float64 `json:"stake"`
	PotentialWin float64 `json:"potential_win"`
	CID          string  `json:"cid"`
}
