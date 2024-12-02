package types

type Team struct {
	ID     int    `bun:"id" json:"id"`
	Name   string `bun:"name" json:"name"`
	Logo   string `bun:"logo" json:"logo"`
	Winner bool   `bun:"winner" json:"winner"`
}

type Teams struct {
	Home Team `bun:"home" json:"home"`
	Away Team `bun:"away" json:"away"`
}
