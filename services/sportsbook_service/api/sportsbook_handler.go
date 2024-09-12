package api

import (
	"betcube_engine/types"
	"betcube_engine/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"

	datastore "betcube_engine/datastore"
	errors "betcube_engine/errors"
)

type SportsbookHandler struct {
	store *datastore.MongoStore
}

func NewSportsbookHandler(s *datastore.MongoStore) *SportsbookHandler {
	return &SportsbookHandler{
		store: s,
	}
}

func (h *SportsbookHandler) HandlePostSportsbook(ctx *fiber.Ctx) error {
	var sportsbook *types.Sportsbook
	if err := ctx.BodyParser(&sportsbook); err != nil {
		return errors.ErrBadRequest()
	}
	insertedSportsbook, err := h.store.Sportsbook.InsertSportsbook(ctx.Context(), sportsbook)
	if err != nil {
		return errors.ErrBadRequest()
	}
	return ctx.JSON(insertedSportsbook)
}

func (h *SportsbookHandler) HandleGetSportsbookByDate(ctx *fiber.Ctx) error {
	dateStr := ctx.Params("date")
	timezone, err := utils.GetClientTimezone()
	if err != nil {
		return errors.ErrInternalServer()
	}
	sportsbookFixtures, err := h.store.Sportsbook.GetSportsbookByDate(ctx.Context(), dateStr, timezone)
	if err != nil {
		return errors.ErrNotResourceNotFound("with specific date")
	}
	return ctx.JSON(sportsbookFixtures)
}

func (h *SportsbookHandler) HandlerGetSportsbookByFixtureID(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return err
	}
	sportsbook, err := h.store.Sportsbook.GetSportsbookByFixtureID(ctx.Context(), id)
	if err != nil {
		return errors.ErrNotResourceNotFound("sportsbook")
	}
	return ctx.JSON(sportsbook)
}
