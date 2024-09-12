package api

import (
	"betcube_engine/types"

	"github.com/gofiber/fiber/v2"

	datastore "betcube_engine/datastore"
	errors "betcube_engine/errors"
)

type RecordHandler struct {
	store *datastore.SupabaseStore
}

func NewRecordHandler(s *datastore.SupabaseStore) *RecordHandler {
	return &RecordHandler{
		store: s,
	}
}

func (h *RecordHandler) HandlePostRecord(c *fiber.Ctx) error {
	var record *types.Record
	if err := c.BodyParser(&record); err != nil {
		return err
	}
	insertedRecord, err := h.store.Record.InsertRecord(c.Context(), record)
	if err != nil {
		return errors.ErrBadRequest()
	}
	return c.JSON(insertedRecord)
}

func (h *RecordHandler) HandleGetRecordsByUser(c *fiber.Ctx) error {
	user := c.Params("wallet_address")
	records, err := h.store.Record.GetRecordsByUser(c.Context(), user)
	if err != nil {
		return errors.ErrNotResourceNotFound("wallet_address")
	}
	return c.JSON(records)
}
