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
	"time"

	"github.com/joho/godotenv"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// API Response structures
type APIResponse struct {
	Get        string       `json:"get"`
	Parameters Parameters   `json:"parameters"`
	Results    int          `json:"results"`
	Response   []APIFixture `json:"response"`
}

type Parameters struct {
	Date string `json:"date"`
}

// API Models
type APIFixture struct {
	Fixture struct {
		ID     int       `json:"id"`
		Date   string    `json:"date"`
		Status APIStatus `json:"status"`
	} `json:"fixture"`
	League APILeague `json:"league"`
	Teams  APITeams  `json:"teams"`
	Goals  APIGoals  `json:"goals"`
	Score  APIScore  `json:"score"`
}

type APIStatus struct {
	Long    string `json:"long"`
	Short   string `json:"short"`
	Elapsed int    `json:"elapsed"`
}

type APILeague struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
	Logo    string `json:"logo"`
	Flag    string `json:"flag"`
}

type APITeams struct {
	Home APITeam `json:"home"`
	Away APITeam `json:"away"`
}

type APITeam struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Logo   string `json:"logo"`
	Winner bool   `json:"winner"`
}

type APIGoals struct {
	Home int `json:"home"`
	Away int `json:"away"`
}

type APIScore struct {
	Halftime  APIScoreDetail `json:"halftime"`
	Fulltime  APIScoreDetail `json:"fulltime"`
	Extratime APIScoreDetail `json:"extratime"`
	Penalty   APIScoreDetail `json:"penalty"`
}

type APIScoreDetail struct {
	Home *int `json:"home"`
	Away *int `json:"away"`
}

// Database Models
type League struct {
	bun.BaseModel `bun:"table:leagues,alias:l"`

	ID      int    `bun:"id,pk"`
	Name    string `bun:"name,notnull"`
	Country string `bun:"country"`
	Logo    string `bun:"logo"`
	Flag    string `bun:"flag"`
}

type Team struct {
	bun.BaseModel `bun:"table:teams,alias:t"`

	ID     int    `bun:"id,pk"`
	Name   string `bun:"name,notnull"`
	Logo   string `bun:"logo"`
	Winner bool   `bun:"winner"`
}

type Fixture struct {
	bun.BaseModel `bun:"table:fixtures,alias:f"`

	ID         int             `bun:"id,pk"`
	Date       time.Time       `bun:"date"`
	LeagueID   int             `bun:"league_id"`
	HomeTeamID int             `bun:"home_team_id"`
	AwayTeamID int             `bun:"away_team_id"`
	Status     json.RawMessage `bun:"status,type:jsonb"`
	Score      json.RawMessage `bun:"score,type:jsonb"`
	Goals      json.RawMessage `bun:"goals,type:jsonb"`
	CreatedAt  time.Time       `bun:"created_at"`
	UpdatedAt  time.Time       `bun:"updated_at"`
}

type Store struct {
	db *bun.DB
}

func NewStore(db *bun.DB) *Store {
	return &Store{db: db}
}

func (s *Store) InsertLeague(ctx context.Context, league *League) error {
	_, err := s.db.NewInsert().
		Model(league).
		On("CONFLICT (id) DO UPDATE").
		Set("name = EXCLUDED.name").
		Set("country = EXCLUDED.country").
		Set("logo = EXCLUDED.logo").
		Set("flag = EXCLUDED.flag").
		Exec(ctx)
	return err
}

func (s *Store) InsertTeam(ctx context.Context, team *Team) error {
	_, err := s.db.NewInsert().
		Model(team).
		On("CONFLICT (id) DO UPDATE").
		Set("name = EXCLUDED.name").
		Set("logo = EXCLUDED.logo").
		Set("winner = EXCLUDED.winner").
		Exec(ctx)
	return err
}

func (s *Store) InsertFixture(ctx context.Context, fixture *Fixture) error {
	_, err := s.db.NewInsert().
		Model(fixture).
		On("CONFLICT (id) DO UPDATE").
		Set("date = EXCLUDED.date").
		Set("league_id = EXCLUDED.league_id").
		Set("home_team_id = EXCLUDED.home_team_id").
		Set("away_team_id = EXCLUDED.away_team_id").
		Set("status = EXCLUDED.status").
		Set("score = EXCLUDED.score").
		Set("goals = EXCLUDED.goals").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx)
	return err
}

func setupDatabase() (*bun.DB, error) {
	supabaseEndpoint := os.Getenv("SUPABASE_URL")

	connector := pgdriver.NewConnector(
		pgdriver.WithDSN(supabaseEndpoint),
		pgdriver.WithTimeout(10*time.Second),
		pgdriver.WithDialTimeout(5*time.Second),
	)

	sqldb := sql.OpenDB(connector)
	sqldb.SetMaxOpenConns(10)
	sqldb.SetMaxIdleConns(5)
	sqldb.SetConnMaxLifetime(5 * time.Minute)

	db := bun.NewDB(sqldb, pgdialect.New())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("cannot connect to database: %v", err)
	}

	return db, nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	start := time.Now()

	db, err := setupDatabase()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := NewStore(db)
	ctx := context.Background()

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	fixturesEndpoint := "https://v3.football.api-sports.io/fixtures"
	date := "2024-12-03"

	req, err := http.NewRequest("GET", fmt.Sprintf("%s?date=%s", fixturesEndpoint, date), nil)
	if err != nil {
		log.Fatal("Error creating request:", err)
	}

	req.Header.Add("x-rapidapi-host", "v3.football.api-sports.io")
	req.Header.Add("x-rapidapi-key", "33db0d9d9531386613988c43458700a6")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error making request:", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response:", err)
	}

	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		log.Fatal("Error parsing response:", err)
	}

	log.Printf("Found %d fixtures to process", len(apiResp.Response))

	// Process leagues
	leagues := make(map[int]*League)
	for _, fixture := range apiResp.Response {
		if _, exists := leagues[fixture.League.ID]; !exists {
			leagues[fixture.League.ID] = &League{
				ID:      fixture.League.ID,
				Name:    fixture.League.Name,
				Country: fixture.League.Country,
				Logo:    fixture.League.Logo,
				Flag:    fixture.League.Flag,
			}
		}
	}

	log.Printf("Inserting %d leagues...", len(leagues))
	for _, league := range leagues {
		if err := store.InsertLeague(ctx, league); err != nil {
			log.Printf("Warning - league insert: %v", err)
		}
	}

	// Process teams
	teams := make(map[int]*Team)
	for _, fixture := range apiResp.Response {
		if _, exists := teams[fixture.Teams.Home.ID]; !exists {
			teams[fixture.Teams.Home.ID] = &Team{
				ID:     fixture.Teams.Home.ID,
				Name:   fixture.Teams.Home.Name,
				Logo:   fixture.Teams.Home.Logo,
				Winner: fixture.Teams.Home.Winner,
			}
		}
		if _, exists := teams[fixture.Teams.Away.ID]; !exists {
			teams[fixture.Teams.Away.ID] = &Team{
				ID:     fixture.Teams.Away.ID,
				Name:   fixture.Teams.Away.Name,
				Logo:   fixture.Teams.Away.Logo,
				Winner: fixture.Teams.Away.Winner,
			}
		}
	}

	log.Printf("Inserting %d teams...", len(teams))
	for _, team := range teams {
		if err := store.InsertTeam(ctx, team); err != nil {
			log.Printf("Warning - team insert: %v", err)
		}
	}

	// Process fixtures
	successCount := 0
	errorCount := 0

	log.Printf("Inserting fixtures...")
	for i, apiFixture := range apiResp.Response {
		status, _ := json.Marshal(apiFixture.Fixture.Status)
		score, _ := json.Marshal(apiFixture.Score)
		goals, _ := json.Marshal(apiFixture.Goals)

		date, err := time.Parse("2006-01-02T15:04:05-07:00", apiFixture.Fixture.Date)
		if err != nil {
			log.Printf("Warning - date parse error for fixture %d: %v", apiFixture.Fixture.ID, err)
			continue
		}

		fixture := &Fixture{
			ID:         apiFixture.Fixture.ID,
			Date:       date,
			LeagueID:   apiFixture.League.ID,
			HomeTeamID: apiFixture.Teams.Home.ID,
			AwayTeamID: apiFixture.Teams.Away.ID,
			Status:     status,
			Score:      score,
			Goals:      goals,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		if err := store.InsertFixture(ctx, fixture); err != nil {
			log.Printf("Error processing fixture %d: %v", apiFixture.Fixture.ID, err)
			errorCount++
		} else {
			successCount++
		}

		if i > 0 && i%50 == 0 {
			log.Printf("Processed %d/%d fixtures...", i, len(apiResp.Response))
		}
	}

	duration := time.Since(start)
	log.Printf("Seeding completed in %s. Success: %d, Errors: %d", duration, successCount, errorCount)
}
