package types

import "time"

type Nfticket struct {
	NFTIdentifier string `bun:"nft_identifier,pk" json:"nft_identifier"`
	Collection    string `bun:"collection" json:"collection"`
	WalletAddress string `bun:"wallet_address" json:"wallet_address"`
	// Bets          []Bets    `bun:"bets" json:"bets"`
	TotalOdds    float64   `bun:"total_odds" json:"total_odds"`
	Stake        float64   `bun:"stake" json:"stake"`
	PotentialWin float64   `bun:"potential_win" json:"potential_win"`
	Minted       time.Time `bun:"minted" json:"minted"`
	Result       string    `bun:"result" json:"result"`
	IsPaid       bool      `bun:"is_paid" json:"is_paid"`
	Explorer     string    `bun:"explorer" json:"explorer"`
	Nonce        int       `bun:"nonce" json:"nonce"`
}
