package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"rockstake-core/types"

	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// API Response structures
type APIResponse struct {
	Get        string          `json:"get"`
	Parameters Parameters      `json:"parameters"`
	Results    int             `json:"results"`
	Response   []types.Fixture `json:"response"`
}

type Parameters struct {
	Date string `json:"date"`
}

type Store struct {
	db *bun.DB
}

func NewStore(db *bun.DB) *Store {
	return &Store{db: db}
}

func (s *Store) InsertLeague(ctx context.Context, league *types.League) error {
	_, err := s.db.NewInsert().Model(league).Exec(ctx)
	return err
}

func (s *Store) InsertTeam(ctx context.Context, team *types.Team) error {
	_, err := s.db.NewInsert().Model(team).Exec(ctx)
	return err
}

func (s *Store) InsertFixture(ctx context.Context, fixture *types.Fixture) error {
	_, err := s.db.NewInsert().Model(fixture).Exec(ctx)
	return err
}

func main() {
	ctx := context.Background()

	supabaseEndpoint := os.Getenv("SUPABASE_URL")
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(supabaseEndpoint)))
	dbSupabase := bun.NewDB(sqldb, pgdialect.New())

	defer dbSupabase.Close()
	store := NewStore(dbSupabase)

	fixturesEndpoint := "https://v3.football.api-sports.io/fixtures"
	date := "2024-12-02"
	headers := map[string]string{
		"x-rapidapi-host": "v3.football.api-sports.io",
		"x-rapidapi-key":  "33db0d9d9531386613988c43458700a6",
	}

	// Create HTTP request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?date=%s", fixturesEndpoint, date), nil)
	if err != nil {
		log.Fatal(err)
	}

	for key, value := range headers {
		req.Header.Add(key, value)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Parse response
	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		log.Fatal(err)
	}

	// Process each fixture
	for _, fixture := range apiResp.Response {
		// Insert league
		league := &types.League{
			ID:      fixture.League.ID,
			Name:    fixture.League.Name,
			Country: fixture.League.Country,
			Logo:    fixture.League.Logo,
			Flag:    fixture.League.Flag,
		}
		if err := store.InsertLeague(ctx, league); err != nil {
			log.Printf("Error inserting league: %v", err)
			continue
		}

		// Insert teams
		homeTeam := &types.Team{
			ID:   fixture.Teams.Home.ID,
			Name: fixture.Teams.Home.Name,
			Logo: fixture.Teams.Home.Logo,
		}
		if err := store.InsertTeam(ctx, homeTeam); err != nil {
			log.Printf("Error inserting home team: %v", err)
			continue
		}

		awayTeam := &types.Team{
			ID:   fixture.Teams.Away.ID,
			Name: fixture.Teams.Away.Name,
			Logo: fixture.Teams.Away.Logo,
		}
		if err := store.InsertTeam(ctx, awayTeam); err != nil {
			log.Printf("Error inserting away team: %v", err)
			continue
		}

		// Insert fixture
		fixtureDB := &types.Fixture{
			ID:     fixture.ID,
			Date:   fixture.Date,
			League: fixture.League,
			Teams:  fixture.Teams,
			Status: fixture.Status,
			Score:  fixture.Score,
			Goals:  fixture.Goals,
		}
		if err := store.InsertFixture(ctx, fixtureDB); err != nil {
			log.Printf("Error inserting fixture: %v", err)
			continue
		}
	}

	log.Println("Data seeding completed successfully")
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
