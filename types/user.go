package types

type User struct {
	WalletAddress string `bun:"wallet_address" json:"wallet_address"`
	ClientId      string `bun:"user_id" json:"user_id"`
	ShareSeed     string `bun:"share_seed" json:"share_seed"`
}
