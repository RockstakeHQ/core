package types

import "time"

type Sportsbook struct {
	ID     int       `bson:"_id" json:"_id"`
	Update time.Time `bson:"update" json:"update"`
	Bets   []Bet     `bson:"bets" json:"bets"`
}

type Bet struct {
	ID     int     `bson:"id" json:"id"`
	Name   string  `bson:"name" json:"name"`
	Values []Value `bson:"values" json:"values"`
}

type Value struct {
	ID    int    `bson:"id" json:"id"`
	Value string `bson:"value" json:"value"`
	Odd   string `bson:"odd" json:"odd"`
}
