package db

import (
	"context"
	"fmt"
	"log"
	"rockstake-core/types"

	"github.com/pocketbase/pocketbase"
)

type FixtureStore interface {
	GetFixturesByDateRange(ctx context.Context, startDate, endDate string) ([]types.Fixture, error)
}

type PocketBaseFixtureStore struct {
	app *pocketbase.PocketBase
}

func NewPocketBaseFixtureStore(app *pocketbase.PocketBase) *PocketBaseFixtureStore {
	return &PocketBaseFixtureStore{
		app: app,
	}
}

func (s *PocketBaseFixtureStore) GetFixturesByDateRange(ctx context.Context, startDate, endDate string) ([]types.Fixture, error) {
	log.Printf("Searching fixtures between %s and %s", startDate, endDate)

	query := `
        SELECT 
            id,           
            fixture_id, 
            date,
            venue_name,
            venue_city,
            league_name,
            league_country,
            home_team,
            away_team,
            status,
            home_goals,
            away_goals
        FROM fixtures 
        WHERE date >= {:start_date} AND date < {:end_date}
    `

	var fixtures []types.Fixture
	err := s.app.DB().NewQuery(query).
		Bind(map[string]interface{}{
			"start_date": startDate,
			"end_date":   endDate,
		}).
		All(&fixtures)

	if err != nil {
		return nil, fmt.Errorf("error fetching fixtures: %v", err)
	}

	for _, f := range fixtures {
		log.Printf("Found fixture: ID=%d, League=%s/%s, Teams=%s vs %s",
			f.ID, f.LeagueName, f.LeagueCountry, f.HomeTeam, f.AwayTeam)
	}

	return fixtures, nil
}
