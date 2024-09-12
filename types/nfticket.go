package types

import "time"

type Nfticket struct {
	NFTIdentifier string    `bun:"nft_identifier,pk" json:"nft_identifier"`
	Collection    string    `bun:"collection" json:"collection"`
	WalletAddress string    `bun:"wallet_address" json:"wallet_address"`
	Bets          []Bets    `bun:"bets" json:"bets"`
	TotalOdds     float64   `bun:"total_odds" json:"total_odds"`
	Stake         float64   `bun:"stake" json:"stake"`
	PotentialWin  float64   `bun:"potential_win" json:"potential_win"`
	Minted        time.Time `bun:"minted" json:"minted"`
	Result        string    `bun:"result" json:"result"`
	IsPaid        bool      `bun:"is_paid" json:"is_paid"`
	XSpotlight    string    `bun:"xspotlight" json:"xspotlight"`
	Explorer      string    `bun:"explorer" json:"explorer"`
	Nonce         int       `bun:"nonce" json:"nonce"`
}

type Bets struct {
	FixtureID  int     `bun:"fixture_id" json:"fixture_id"`
	BetTypeID  int     `bun:"bet_type" json:"bet_type"`
	BetValueID int     `bun:"bet_value" json:"bet_value"`
	Odd        float64 `bun:"odd" json:"odd"`
	Result     string  `bun:"result" json:"result"`
}

type EnrichedNfticket struct {
	NFTIdentifier string        `bun:"nft_identifier,pk" json:"nft_identifier"`
	Collection    string        `bun:"collection" json:"collection"`
	WalletAddress string        `bun:"wallet_address" json:"wallet_address"`
	EnrichedBets  []EnrichedBet `bun:"bets" json:"bets"`
	TotalOdds     float64       `bun:"total_odds" json:"total_odds"`
	Stake         float64       `bun:"stake" json:"stake"`
	PotentialWin  float64       `bun:"potential_win" json:"potential_win"`
	Minted        time.Time     `bun:"minted" json:"minted"`
	Result        string        `bun:"result" json:"result"`
	IsPaid        bool          `bun:"is_paid" json:"is_paid"`
	XSpotlight    string        `bun:"xspotlight" json:"xspotlight"`
	Explorer      string        `bun:"explorer" json:"explorer"`
	Nonce         int           `bun:"nonce" json:"nonce"`
}

type EnrichedBet struct {
	Fixture  FixtureWithoutSportsbook `bun:"fixture" json:"fixture"`
	BetType  BetType                  `bun:"bet_type" json:"bet_type"`
	BetValue BetValue                 `bun:"bet_value" json:"bet_value"`
	Odd      float64                  `bun:"odd" json:"odd"`
	Result   string                   `bun:"result" json:"result"`
}

type UpdatedNfticket struct {
	Bets         []Bets  `bun:"bets" json:"bets"`
	TotalOdds    float64 `bun:"total_odds" json:"total_odds"`
	PotentialWin float64 `bun:"potential_win" json:"potential_win"`
	Result       string  `bun:"result" json:"result"`
	IsPaid       bool    `bun:"is_paid" json:"is_paid"`
}

type UpdatedNfticketByUser struct {
	IsPaid bool `bun:"is_paid" json:"is_paid"`
}
