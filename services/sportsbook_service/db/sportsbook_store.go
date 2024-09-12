package db

import (
	"betcube_engine/types"
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	utils "betcube_engine/utils"
)

type SportsbookStore interface {
	InsertSportsbook(context.Context, *types.Sportsbook) (*types.Sportsbook, error)
	GetSportsbookByDate(context.Context, string, string) ([]*types.QueryWinnerBetsFixture, error)
	GetSportsbookByFixtureID(context.Context, int) (*types.QueryWinnerBetsFixture, error)
}

type MongoSportsbookStore struct {
	client                *mongo.Client
	sportsbooksCollection *mongo.Collection
	fixturesCollection    *mongo.Collection
	leaguesCollection     *mongo.Collection
	teamsCollection       *mongo.Collection
}

func NewMongoSportsbookStore(client *mongo.Client) *MongoSportsbookStore {
	dbName := os.Getenv("MONGO_DB_NAME")
	return &MongoSportsbookStore{
		client:                client,
		sportsbooksCollection: client.Database(dbName).Collection("sportsbooks"),
		fixturesCollection:    client.Database(dbName).Collection("fixtures"),
		leaguesCollection:     client.Database(dbName).Collection("leagues"),
		teamsCollection:       client.Database(dbName).Collection("teams"),
	}
}

func (s *MongoSportsbookStore) InsertSportsbook(ctx context.Context, sportsbook *types.Sportsbook) (*types.Sportsbook, error) {
	_, err := s.sportsbooksCollection.InsertOne(ctx, sportsbook)
	if err != nil {
		return nil, err
	}
	return sportsbook, nil
}

func (s *MongoSportsbookStore) GetSportsbookByFixtureID(ctx context.Context, id int) (*types.QueryWinnerBetsFixture, error) {
	var fixture types.Fixture
	if err := s.fixturesCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&fixture); err != nil {
		return nil, err
	}

	var league types.League
	err := s.leaguesCollection.FindOne(ctx, bson.M{"_id": fixture.League}).Decode(&league)
	if err != nil {
		return nil, err
	}
	var homeTeam, awayTeam types.Team
	err = s.teamsCollection.FindOne(ctx, bson.M{"_id": fixture.Teams[0]}).Decode(&homeTeam)
	if err != nil {
		return nil, err
	}
	err = s.teamsCollection.FindOne(ctx, bson.M{"_id": fixture.Teams[1]}).Decode(&awayTeam)
	if err != nil {
		return nil, err
	}

	var sportsbook types.Sportsbook
	if err := s.sportsbooksCollection.FindOne(ctx, bson.M{"_id": fixture.SportsbookID}).Decode(&sportsbook); err != nil {
		return nil, err
	}

	var queryResult = types.QueryWinnerBetsFixture{
		ID:         fixture.ID,
		Date:       fixture.Date,
		League:     league,
		Teams:      types.Teams{Home: homeTeam, Away: awayTeam},
		Status:     fixture.Status,
		Score:      fixture.Score,
		Goals:      fixture.Goals,
		Sportsbook: sportsbook.Bets,
	}
	return &queryResult, nil
}

func (s *MongoSportsbookStore) GetSportsbookByDate(ctx context.Context, dateString string, timezone string) ([]*types.QueryWinnerBetsFixture, error) {
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return nil, err
	}
	adjustedDate, err := utils.AdjustTimeToUserTimezone(date, timezone)
	if err != nil {
		return nil, err
	}
	endDate := adjustedDate.Add(24 * time.Hour)

	filter := bson.M{
		"date": bson.M{
			"$gte": adjustedDate,
			"$lt":  endDate,
		},
	}

	cur, err := s.fixturesCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var queryResults []*types.QueryWinnerBetsFixture
	for cur.Next(ctx) {
		var fixture types.Fixture
		if err := cur.Decode(&fixture); err != nil {
			return nil, err
		}

		adjustedTime, err := utils.AdjustTimeToUserTimezone(fixture.Date, timezone)
		if err != nil {
			return nil, err
		}
		fixture.Date = adjustedTime

		currentTime := time.Now().UTC()
		currentDate := currentTime.Truncate(12 * time.Hour)
		fixtureDate := adjustedDate.Truncate(24 * time.Hour)

		// Verificăm dacă data selectată este ziua curentă
		if fixtureDate.Equal(currentDate) {
			// Dacă este ziua curentă, excludem meciurile care încep într-un minut
			if fixture.Date.Before(currentTime.Add(1 * time.Minute)) {
				continue
			}
		}

		var league types.League
		err = s.leaguesCollection.FindOne(ctx, bson.M{"_id": fixture.League}).Decode(&league)
		if err != nil {
			return nil, err
		}

		var homeTeam, awayTeam types.Team
		err = s.teamsCollection.FindOne(ctx, bson.M{"_id": fixture.Teams[0]}).Decode(&homeTeam)
		if err != nil {
			return nil, err
		}
		err = s.teamsCollection.FindOne(ctx, bson.M{"_id": fixture.Teams[1]}).Decode(&awayTeam)
		if err != nil {
			return nil, err
		}

		var sportsbook types.Sportsbook
		if err := s.sportsbooksCollection.FindOne(ctx, bson.M{"_id": fixture.SportsbookID}).Decode(&sportsbook); err != nil {
			return nil, err
		}

		winnerBets := filterWinnerBets(sportsbook.Bets)

		queryResult := &types.QueryWinnerBetsFixture{
			ID:         fixture.ID,
			Date:       fixture.Date,
			League:     league,
			Teams:      types.Teams{Home: homeTeam, Away: awayTeam},
			Sportsbook: winnerBets,
			Status:     fixture.Status,
			Score:      fixture.Score,
			Goals:      fixture.Goals,
		}

		queryResults = append(queryResults, queryResult)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}
	if len(queryResults) == 0 {
		return nil, fmt.Errorf("no fixtures found for the specified date")
	}
	return queryResults, nil
}

func filterWinnerBets(bets []types.Bet) []types.Bet {
	var winnerBets []types.Bet
	for _, bet := range bets {
		if bet.ID == 1 {
			winnerBets = append(winnerBets, bet)
		}
	}
	return winnerBets
}
