package types

import "time"

type Fixture struct {
	ID           int       `bson:"_id" json:"id"`
	Date         time.Time `bson:"date" json:"date"`
	League       int       `bson:"league" json:"league"`
	Teams        []int     `bson:"teams" json:"teams"`
	Status       Status    `bson:"status,omitempty" json:"status,omitempty"`
	Score        Score     `bson:"score,omitempty" json:"score,omitempty"`
	Goals        Goals     `bson:"goals,omitempty" json:"goals,omitempty"`
	SportsbookID int       `bson:"sportsbook_id" json:"sportsbook_id"`
}

type Status struct {
	Long    string `bson:"long" json:"long"`
	Short   string `bson:"short" json:"short"`
	Elapsed int    `bson:"elapsed" json:"elapsed"`
}

type FixtureWithoutSportsbook struct {
	ID     int       `bson:"_id" json:"id"`
	Date   time.Time `bson:"date" json:"date"`
	League League    `bson:"league" json:"league"`
	Teams  Teams     `bson:"teams" json:"teams"`
	Score  Score     `bson:"score" json:"score"`
	Goals  Goals     `bson:"goals" json:"goals"`
	Status Status    `bson:"status" json:"status"`
}

type QueryWinnerBetsFixture struct {
	ID         int       `bson:"_id" json:"id"`
	Date       time.Time `bson:"date" json:"date"`
	League     League    `bson:"league" json:"league"`
	Teams      Teams     `bson:"teams" json:"teams"`
	Sportsbook []Bet     `bson:"bets" json:"bets"`
	Status     Status    `bson:"status,omitempty" json:"status,omitempty"`
	Score      Score     `bson:"score,omitempty" json:"score,omitempty"`
	Goals      Goals     `bson:"goals,omitempty" json:"goals,omitempty"`
}

type QuerySportsbookFixture struct {
	ID         int        `bson:"_id" json:"id"`
	Date       time.Time  `bson:"date" json:"date"`
	League     League     `bson:"league" json:"league"`
	Teams      Teams      `bson:"teams" json:"teams"`
	Status     Status     `bson:"status" json:"status"`
	Score      Score      `bson:"score" json:"score"`
	Goals      Goals      `bson:"goals" json:"goals"`
	Sportsbook Sportsbook `bson:"sportbook" json:"sportbook"`
}
