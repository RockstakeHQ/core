package types

import "time"

type Fixture struct {
	ID     int       `bun:"id" json:"id"`
	Date   time.Time `bun:"date" json:"date"`
	League League    `bun:"league" json:"league"`
	Teams  Teams     `bun:"teams" json:"teams"`
	Status Status    `bun:"status" json:"status"`
	Score  Score     `bun:"score" json:"score"`
	Goals  Goals     `bun:"goals" json:"goals"`
}

type Status struct {
	Long    string `bun:"long" json:"long"`
	Short   string `bun:"short" json:"short"`
	Elapsed int    `bun:"elapsed" json:"elapsed"`
}
