package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"rockstake-core/types"

	"github.com/pocketbase/pocketbase"
)

type BetStore interface {
	InsertBet(ctx context.Context, bet types.Bet) (types.Bet, error)
	GetBet(ctx context.Context, nft_identifier string) (types.Bet, error)
}

type PocketBaseBetStore struct {
	app *pocketbase.PocketBase
}

func NewPocketBaseBetStore(app *pocketbase.PocketBase) *PocketBaseBetStore {
	return &PocketBaseBetStore{
		app: app,
	}
}

func (s *PocketBaseBetStore) InsertBet(ctx context.Context, bet types.Bet) (types.Bet, error) {
	betInfoJSON, err := json.Marshal(bet.BetInfo)
	if err != nil {
		return types.Bet{}, fmt.Errorf("error marshaling bet info: %v", err)
	}

	query := `
        INSERT INTO bets (
			nft_identifier,
			collection,
			nonce
            wallet_address,
            fixture_id,
            market_id,
            status,
            bet_info
        ) VALUES (
		    {:nft_identifier},
			{:collection},
			{:nonce},
            {:wallet_address},
            {:fixture_id},
            {:market_id},
            {:status},
            {:bet_info}
        )
    `

	params := map[string]interface{}{
		"nft_identifier": bet.NftIdentifier,
		"collection":     bet.Collection,
		"nonce":          bet.Nonce,
		"wallet_address": bet.WalletAddress,
		"fixture_id":     bet.FixtureID,
		"market_id":      bet.MarketID,
		"status":         bet.Status,
		"bet_info":       string(betInfoJSON),
	}

	_, err = s.app.DB().NewQuery(query).
		Bind(params).
		Execute()

	if err != nil {
		return types.Bet{}, fmt.Errorf("error inserting bet: %v", err)
	}

	log.Printf("Inserted bet: Wallet=%s, Fixture=%d, Market=%d, Status=%s",
		bet.WalletAddress, bet.FixtureID, bet.MarketID, bet.Status)

	return bet, nil
}

func (s *PocketBaseBetStore) GetBet(ctx context.Context, nft_identifier string) (types.Bet, error) {
	query := `
        SELECT *   
        FROM bets 
        WHERE id = {:nft_identifier}
    `

	var bet types.Bet
	err := s.app.DB().NewQuery(query).
		Bind(map[string]interface{}{
			"nft_identifier": nft_identifier,
		}).
		One(&bet)

	if err != nil {
		return types.Bet{}, fmt.Errorf("error fetching bet: %w", err)
	}

	return bet, nil
}
