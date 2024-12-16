package types

type Fixture struct {
	ID            string `json:"id"`
	FixtureID     int    `json:"fixture_id"`
	Date          string `json:"date"`
	VenueName     string `json:"venue_name"`
	VenueCity     string `json:"venue_city"`
	LeagueName    string `json:"league_name"`
	LeagueCountry string `json:"league_country"`
	HomeTeam      string `json:"home_team"`
	AwayTeam      string `json:"away_team"`
	Status        string `json:"status"`
	HomeGoals     int    `json:"home_goals"`
	AwayGoals     int    `json:"away_goals"`
}
