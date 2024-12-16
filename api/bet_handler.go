package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rockstake-core/db"
	"rockstake-core/types"

	"github.com/pocketbase/pocketbase/core"
)

type BetHandler struct {
	store db.BetStore
}

func NewBetHandler(store db.BetStore) *BetHandler {
	return &BetHandler{
		store: store,
	}
}

func (h *BetHandler) HandlePostBet(e *core.RequestEvent) error {
	var betRequest struct {
		WalletAddress string  `json:"wallet_address"`
		FixtureID     int     `json:"fixture_id"`
		MarketID      int     `json:"market_id"`
		Selection     string  `json:"selection"`
		Type          string  `json:"type"`
		Odd           float64 `json:"odd"`
		Stake         float64 `json:"stake"`
		CID           string  `json:"cid"`
	}

	if err := json.NewDecoder(e.Request.Body).Decode(&betRequest); err != nil {
		http.Error(e.Response, fmt.Sprintf("Error parsing request: %v", err), http.StatusBadRequest)
		return nil
	}

	if betRequest.WalletAddress == "" || betRequest.FixtureID == 0 || betRequest.MarketID == 0 {
		http.Error(e.Response, "Missing required fields", http.StatusBadRequest)
		return nil
	}

	bet := types.Bet{
		WalletAddress: betRequest.WalletAddress,
		FixtureID:     betRequest.FixtureID,
		MarketID:      betRequest.MarketID,
		Status:        "open",
		BetInfo: types.BetInfo{
			Selection:    betRequest.Selection,
			Type:         betRequest.Type,
			Odd:          betRequest.Odd,
			Stake:        betRequest.Stake,
			PotentialWin: betRequest.Odd * betRequest.Stake,
			CID:          betRequest.CID,
		},
	}

	newBet, err := h.store.InsertBet(e.Request.Context(), bet)
	if err != nil {
		http.Error(e.Response, fmt.Sprintf("Error creating bet: %v", err), http.StatusInternalServerError)
		return nil
	}

	e.Response.Header().Set("Content-Type", "application/json")
	e.Response.WriteHeader(http.StatusCreated)
	return json.NewEncoder(e.Response).Encode(map[string]interface{}{
		"message": "Bet created successfully",
		"bet":     newBet,
	})
}