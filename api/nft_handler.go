package api

import (
	"encoding/json"
	"net/http"
	"rockstake-core/db"
	"rockstake-core/types"

	"github.com/pocketbase/pocketbase/core"
)

type NFTHandler struct {
	store db.NFTStore
}

func NewNFTHandler(store db.NFTStore) *NFTHandler {
	return &NFTHandler{
		store: store,
	}
}

func (h *NFTHandler) HandleGenerateNFT(e *core.RequestEvent) error {
	// Parsăm request-ul
	var betData types.NftNodeInfo
	if err := json.NewDecoder(e.Request.Body).Decode(&betData); err != nil {
		http.Error(e.Response, "Invalid request data", http.StatusBadRequest)
		return nil
	}

	// Validăm datele
	if betData.FixtureID == "" || betData.Market == "" || betData.Selection == "" {
		http.Error(e.Response, "Missing required fields", http.StatusBadRequest)
		return nil
	}

	// Generăm și uploadăm NFT-ul
	cid, err := h.store.GenerateAndUploadNFT(e.Request.Context(), betData)
	if err != nil {
		http.Error(e.Response, err.Error(), http.StatusInternalServerError)
		return nil
	}

	// Returnăm răspunsul
	e.Response.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(e.Response).Encode(map[string]interface{}{
		"message": "NFT metadata generated and uploaded successfully",
		"cid":     cid,
	})
}
