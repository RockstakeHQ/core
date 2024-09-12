package db

import (
	"betcube-engine/types"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/uptrace/bun"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type NfticketStore interface {
	InsertNfticket(context.Context, *types.Nfticket) (*types.Nfticket, error)
	UpdateNfticket(context.Context, string, types.UpdatedNfticket) error
	UpdateNfticketByUser(context.Context, string, types.UpdatedNfticketByUser) error
	GetNfticketsByWalletAddress(context.Context, string, *time.Time, *string) (*[]types.EnrichedNfticket, error)
	GetNfticket(context.Context, string) (*types.EnrichedNfticket, error)
}

type MixedNfticketDBStore struct {
	client              *mongo.Client
	bunClient           *bun.DB
	fixturesCollection  *mongo.Collection
	leaguesCollection   *mongo.Collection
	teamsCollection     *mongo.Collection
	betTypesCollection  *mongo.Collection
	betValuesCollection *mongo.Collection
}

func NewMixedNfticketDBStore(bunClient *bun.DB, mongoClient *mongo.Client) *MixedNfticketDBStore {
	dbName := os.Getenv("MONGO_DB_NAME")
	return &MixedNfticketDBStore{
		bunClient:           bunClient,
		client:              mongoClient,
		fixturesCollection:  mongoClient.Database(dbName).Collection("fixtures"),
		leaguesCollection:   mongoClient.Database(dbName).Collection("leagues"),
		teamsCollection:     mongoClient.Database(dbName).Collection("teams"),
		betTypesCollection:  mongoClient.Database(dbName).Collection("bet_types"),
		betValuesCollection: mongoClient.Database(dbName).Collection("bet_values"),
	}
}

func (s *MixedNfticketDBStore) UpdateNfticket(ctx context.Context, filter string, params types.UpdatedNfticket) error {
	updateStmt := s.bunClient.NewUpdate().Model(&types.Nfticket{}).Where("nft_identifier = ?", filter)
	if params.Result != "" {
		updateStmt.Set("bets = ?", params.Bets)
		updateStmt.Set("total_odds = ?", params.TotalOdds)
		updateStmt.Set("potential_win = ?", params.PotentialWin)
		updateStmt.Set("result = ?", params.Result)
		updateStmt.Set("is_paid = ?", params.IsPaid)
	}
	_, err := updateStmt.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *MixedNfticketDBStore) UpdateNfticketByUser(ctx context.Context, filter string, params types.UpdatedNfticketByUser) error {
	updateStmt := s.bunClient.NewUpdate().Model(&types.Nfticket{}).Where("nft_identifier = ?", filter)
	updateStmt.Set("is_paid = ?", params.IsPaid)
	_, err := updateStmt.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *MixedNfticketDBStore) InsertNfticket(ctx context.Context, nfticket *types.Nfticket) (*types.Nfticket, error) {
	_, err := s.bunClient.NewInsert().Model(nfticket).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return nfticket, nil
}

func (s *MixedNfticketDBStore) GetNfticket(ctx context.Context, nft_identifier string) (*types.EnrichedNfticket, error) {
	var nfticketDatabase types.Nfticket
	query := s.bunClient.NewSelect().Model(&nfticketDatabase).Where("nft_identifier = ?", nft_identifier)

	err := query.Scan(ctx, &nfticketDatabase)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Lol: %+v\n", nfticketDatabase) // Adaugă acest log

	var enrichedNfticket types.EnrichedNfticket

	for _, bet := range nfticketDatabase.Bets {
		fixture, err := s.getFixtureByID(ctx, bet.FixtureID)
		if err != nil {
			return nil, err
		}
		enrichBetType, err := s.getBetType(ctx, bet.BetTypeID)
		if err != nil {
			return nil, err
		}
		enrichBetValue, err := s.getBetValue(ctx, bet.BetValueID)
		if err != nil {
			return nil, err
		}
		var enrichedBet = types.EnrichedBet{
			Fixture:  *fixture,
			BetType:  *enrichBetType,
			BetValue: *enrichBetValue,
			Odd:      bet.Odd,
			Result:   bet.Result,
		}

		enrichedNfticket.EnrichedBets = append(enrichedNfticket.EnrichedBets, enrichedBet)
	}

	enrichedNfticket.NFTIdentifier = nfticketDatabase.NFTIdentifier
	enrichedNfticket.Collection = nfticketDatabase.Collection
	enrichedNfticket.WalletAddress = nfticketDatabase.WalletAddress
	enrichedNfticket.TotalOdds = nfticketDatabase.TotalOdds
	enrichedNfticket.Stake = nfticketDatabase.Stake
	enrichedNfticket.PotentialWin = nfticketDatabase.PotentialWin
	enrichedNfticket.Minted = nfticketDatabase.Minted
	enrichedNfticket.Result = nfticketDatabase.Result
	enrichedNfticket.IsPaid = nfticketDatabase.IsPaid
	enrichedNfticket.XSpotlight = nfticketDatabase.XSpotlight
	enrichedNfticket.Explorer = nfticketDatabase.Explorer
	enrichedNfticket.Nonce = nfticketDatabase.Nonce

	return &enrichedNfticket, nil
}

func (s *MixedNfticketDBStore) GetNfticketsByWalletAddress(ctx context.Context, wallet string, minted *time.Time, result *string) (*[]types.EnrichedNfticket, error) {
	var nftickets []types.Nfticket
	query := s.bunClient.NewSelect().Model(&nftickets).Where("wallet_address = ?", wallet)

	if minted != nil {
		query = query.Where("DATE(minted) = ?", minted.Format("2006-01-02"))
	}
	if result != nil {
		query = query.Where("result = ?", *result)
	}

	err := query.Scan(ctx, &nftickets)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Fetched Nftickets: %+v\n", nftickets) // Adaugă acest log

	var nftickets2 []types.EnrichedNfticket
	for _, nfticket := range nftickets {
		var enrichedNfticket types.EnrichedNfticket
		enrichedNfticket.NFTIdentifier = nfticket.NFTIdentifier
		enrichedNfticket.Collection = nfticket.Collection
		enrichedNfticket.WalletAddress = nfticket.WalletAddress
		enrichedNfticket.TotalOdds = nfticket.TotalOdds
		enrichedNfticket.Stake = nfticket.Stake
		enrichedNfticket.PotentialWin = nfticket.PotentialWin
		enrichedNfticket.Minted = nfticket.Minted
		enrichedNfticket.Result = nfticket.Result
		enrichedNfticket.IsPaid = nfticket.IsPaid
		enrichedNfticket.XSpotlight = nfticket.XSpotlight
		enrichedNfticket.Explorer = nfticket.Explorer
		enrichedNfticket.Nonce = nfticket.Nonce

		for _, bet := range nfticket.Bets {
			fixture, err := s.getFixtureByID(ctx, bet.FixtureID)
			if err != nil {
				return nil, err
			}
			enrichBetType, err := s.getBetType(ctx, bet.BetTypeID)
			if err != nil {
				return nil, err
			}
			enrichBetValue, err := s.getBetValue(ctx, bet.BetValueID)
			if err != nil {
				return nil, err
			}
			var enrichedBet types.EnrichedBet
			enrichedBet.Fixture = *fixture
			enrichedBet.BetType = *enrichBetType
			enrichedBet.BetValue = *enrichBetValue
			enrichedBet.Odd = bet.Odd
			enrichedBet.Result = bet.Result

			enrichedNfticket.EnrichedBets = append(enrichedNfticket.EnrichedBets, enrichedBet)
		}
		nftickets2 = append(nftickets2, enrichedNfticket)
	}

	return &nftickets2, nil
}

func (s *MixedNfticketDBStore) getFixtureByID(ctx context.Context, fixtureID int) (*types.FixtureWithoutSportsbook, error) {
	var fixture types.Fixture
	if err := s.fixturesCollection.FindOne(ctx, bson.M{"_id": fixtureID}).Decode(&fixture); err != nil {
		return nil, err
	}
	var league types.League
	if err := s.leaguesCollection.FindOne(ctx, bson.M{"_id": fixture.League}).Decode(&league); err != nil {
		return nil, err
	}

	var homeTeam, awayTeam types.Team
	err := s.teamsCollection.FindOne(ctx, bson.M{"_id": fixture.Teams[0]}).Decode(&homeTeam)
	if err != nil {
		return nil, err
	}
	err = s.teamsCollection.FindOne(ctx, bson.M{"_id": fixture.Teams[1]}).Decode(&awayTeam)
	if err != nil {
		return nil, err
	}

	var queryResult = types.FixtureWithoutSportsbook{
		ID:     fixture.ID,
		Date:   fixture.Date,
		League: league,
		Teams:  types.Teams{Home: homeTeam, Away: awayTeam},
		Score:  fixture.Score,
		Goals:  fixture.Goals,
		Status: fixture.Status,
	}
	return &queryResult, nil
}

func (s *MixedNfticketDBStore) getBetType(ctx context.Context, betTypeID int) (*types.BetType, error) {
	var betType types.BetType
	if err := s.betTypesCollection.FindOne(ctx, bson.M{"_id": betTypeID}).Decode(&betType); err != nil {
		return nil, err
	}
	return &betType, nil
}

func (s *MixedNfticketDBStore) getBetValue(ctx context.Context, betValueID int) (*types.BetValue, error) {
	var betValue types.BetValue
	if err := s.betValuesCollection.FindOne(ctx, bson.M{"_id": betValueID}).Decode(&betValue); err != nil {
		return nil, err
	}
	return &betValue, nil
}
