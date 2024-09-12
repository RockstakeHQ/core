package types_api

import (
	"betcube-engine/types"
	"encoding/json"
	"fmt"
)

func (v *Value) UnmarshalJSON(data []byte) error {
	var temp interface{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	switch val := temp.(type) {
	case map[string]interface{}:
		if valueStr, ok := val["value"].(string); ok {
			v.Value = valueStr
		} else if valueNum, ok := val["value"].(float64); ok {
			v.Value = fmt.Sprintf("%f", valueNum)
		}
		if oddStr, ok := val["odd"].(string); ok {
			v.Odd = oddStr
		} else if oddNum, ok := val["odd"].(float64); ok {
			v.Odd = fmt.Sprintf("%f", oddNum)
		}
	case string:
		v.Value = val
	case float64:
		v.Value = fmt.Sprintf("%f", val)
	default:
		return fmt.Errorf("unexpected type for Value: %T", val)
	}

	return nil
}

type FixtureResponse struct {
	Response []Fixture `json:"response"`
}

type Fixture struct {
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

type Status struct {
	Long    string `json:"long"`
	Short   string `json:"short"`
	Elapsed int    `json:"elapsed"`
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

type Team struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Logo   string `json:"logo"`
	Winner *bool  `json:"winner"`
}

type GoalsDetails struct {
	Home *int `json:"home"`
	Away *int `json:"away"`
}

type ScoreDetails struct {
	Halftime Score `json:"halftime"`
	Fulltime Score `json:"fulltime"`
}

type Score struct {
	Home *int `json:"home"`
	Away *int `json:"away"`
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
	Bets []Bet `json:"bets"`
}

type Bet struct {
	ID     int     `json:"id"`
	Name   string  `json:"name"`
	Values []Value `json:"values"`
}

type Value struct {
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
	Fixture FixtureAPI   `json:"fixture"`
	League  types.League `json:"league"`
	Teams   types.Teams  `json:"teams"`
	Goals   types.Goals  `json:"goals"`
	Score   types.Score  `json:"score"`
	Update  string       `json:"update"`
	Bets    []types.Bet  `json:"bets"`
}

type FixtureAPI struct {
	ID     int          `json:"id"`
	Date   string       `json:"date"`
	Status types.Status `json:"status"`
}

type MatchDb struct {
	Fixture types.Fixture `json:"fixture"`
	League  types.League  `json:"league"`
	Teams   types.Teams   `json:"teams"`
	Goals   types.Goals   `json:"goals"`
	Score   types.Score   `json:"score"`
	Update  string        `json:"update"`
	Bets    []types.Bet   `json:"bets"`
}
