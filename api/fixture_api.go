package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rockstake-core/db"
	"time"

	"github.com/pocketbase/pocketbase/core"
)

type FixtureHandler struct {
	store db.FixtureStore
}

func NewFixtureHandler(store db.FixtureStore) *FixtureHandler {
	return &FixtureHandler{
		store: store,
	}
}

func (h *FixtureHandler) HandleGetFixturesByDate(e *core.RequestEvent) error {
	// Primim data în formatul "2024-12-16"
	dateStr := e.Request.URL.Query().Get("date")
	if dateStr == "" {
		http.Error(e.Response, "Date parameter is required", http.StatusBadRequest)
		return nil
	}

	// Parsăm data primită
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(e.Response, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return nil
	}

	// Construim intervalul pentru o zi întreagă
	startTime := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endTime := startTime.Add(24 * time.Hour)

	// Formatăm datele pentru query în formatul ISO 8601
	startTimeStr := startTime.Format(time.RFC3339)
	endTimeStr := endTime.Format(time.RFC3339)

	// Folosim intervalul în query
	fixtures, err := h.store.GetFixturesByDateRange(e.Request.Context(), startTimeStr, endTimeStr)
	if err != nil {
		http.Error(e.Response, fmt.Sprintf("Error fetching fixtures: %v", err), http.StatusInternalServerError)
		return nil
	}

	e.Response.Header().Set("Content-Type", "application/json")
	e.Response.WriteHeader(http.StatusOK)

	return json.NewEncoder(e.Response).Encode(map[string]interface{}{
		"fixtures": fixtures,
		"total":    len(fixtures),
	})
}
