package types

type FixtureResponse struct {
	Response []FixtureAPI `json:"response"`
}

type FixtureAPI struct {
	Fixture FixtureDetails `json:"fixture"`
	League  LeagueDetails  `json:"league"`
	Teams   TeamsDetails   `json:"teams"`
	Goals   GoalsDetails   `json:"goals"`
	Score   ScoreDetails   `json:"score"`
}

type FixtureDetails struct {
	ID     int    `json:"id"`
	Date   string `json:"date"`
	Status Status `json:"status"`
}

type LeagueDetails struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
	Logo    string `json:"logo"`
	Flag    string `json:"flag"`
}

type TeamsDetails struct {
	Home Team `json:"home"`
	Away Team `json:"away"`
}

type GoalsDetails struct {
	Home *int `json:"home"`
	Away *int `json:"away"`
}

type ScoreDetails struct {
	Halftime Score `json:"halftime"`
	Fulltime Score `json:"fulltime"`
}

type OddsResponse struct {
	Response []Odds `json:"response"`
	Paging   Paging `json:"paging"`
}

type Odds struct {
	Fixture    FixtureRef  `json:"fixture"`
	Bookmakers []Bookmaker `json:"bookmakers"`
	Update     string      `json:"update"`
}

type FixtureRef struct {
	ID int `json:"id"`
}

type Bookmaker struct {
	Bets []BetAPI `json:"bets"`
}

type BetAPI struct {
	ID     int        `json:"id"`
	Name   string     `json:"name"`
	Values []ValueAPI `json:"values"`
}

type ValueAPI struct {
	Value string `json:"value"`
	Odd   string `json:"odd"`
}

type Paging struct {
	Current int `json:"current"`
	Total   int `json:"total"`
}

type MatchData struct {
	Fixture FixtureDetails `json:"fixture"`
	League  LeagueDetails  `json:"league"`
	Teams   TeamsDetails   `json:"teams"`
	Goals   GoalsDetails   `json:"goals"`
	Score   ScoreDetails   `json:"score"`
	Update  string         `json:"update"`
	Bets    []BetData      `json:"bets"`
}

type BetData struct {
	ID     int         `json:"id"`
	Name   string      `json:"name"`
	Values []ValueData `json:"values"`
}

type ValueData struct {
	Value string `json:"value"`
	Odd   string `json:"odd"`
}

type Match struct {
	Fixture Fixture `json:"fixture"`
	League  League  `json:"league"`
	Teams   Teams   `json:"teams"`
	Goals   Goals   `json:"goals"`
	Score   Score   `json:"score"`
	Update  string  `json:"update"`
	Bets    []Bet   `json:"bets"`
}
