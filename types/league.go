package types

type League struct {
	ID      int    `bun:"id" json:"id"`
	Name    string `bun:"name" json:"name"`
	Country string `bun:"country" json:"country"`
	Logo    string `bun:"logo" json:"logo"`
	Flag    string `bun:"flag" json:"flag"`
}
