package api

import (
	"betcube_engine/types"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"

	datastore "betcube_engine/datastore"
	errors "betcube_engine/errors"
)

type NfticketHandler struct {
	store *datastore.SupabaseStore
}

func NewNfticketHandler(s *datastore.SupabaseStore) *NfticketHandler {
	return &NfticketHandler{
		store: s,
	}
}

func (h *NfticketHandler) HandlePutNfticket(c *fiber.Ctx) error {
	var (
		params         types.UpdatedNfticket
		nft_identifier = c.Params("nft_identifier")
	)
	if err := c.BodyParser(&params); err != nil {
		return errors.ErrBadRequest()
	}
	err := h.store.Nfticket.UpdateNfticket(c.Context(), nft_identifier, params)
	if err != nil {
		return errors.ErrBadRequest()
	}
	return c.JSON(map[string]string{"updated": "ok"})
}

func (h *NfticketHandler) HandlePutNfticketByUser(c *fiber.Ctx) error {
	var (
		params         types.UpdatedNfticketByUser
		nft_identifier = c.Params("nft_identifier")
	)
	if err := c.BodyParser(&params); err != nil {
		return errors.ErrBadRequest()
	}
	err := h.store.Nfticket.UpdateNfticketByUser(c.Context(), nft_identifier, params)
	if err != nil {
		return errors.ErrBadRequest()
	}
	return c.JSON(map[string]string{"updated": "ok"})
}

func (h *NfticketHandler) HandleGetNfticketsByWalletAddress(c *fiber.Ctx) error {
	wallet := c.Params("wallet_address")

	mintedParam := c.Query("minted")
	var minted *time.Time
	if mintedParam != "" {
		parsedMinted, err := time.Parse("2006-01-02", mintedParam)
		if err != nil {
			return errors.ErrBadRequest()
		}
		minted = &parsedMinted
	}

	resultParam := c.Query("result")
	var result *string
	if resultParam != "" {
		result = &resultParam
	}

	nftickets, err := h.store.Nfticket.GetNfticketsByWalletAddress(c.Context(), wallet, minted, result)
	if err != nil {
		return errors.ErrBadRequest()
	}
	return c.JSON(nftickets)
}

func (h *NfticketHandler) HandleGetNfticket(c *fiber.Ctx) error {
	nft_identifier := c.Params("nft_identifier")
	fmt.Println("NFT Identifier:", nft_identifier) // Adaugă acest log

	nfticket, err := h.store.Nfticket.GetNfticket(c.Context(), nft_identifier)
	if err != nil {
		return errors.ErrBadRequest()
	}

	fmt.Printf("Fetched EnrichedNfticket: %+v\n", nfticket) // Adaugă acest log
	return c.JSON(nfticket)
}

func (h *NfticketHandler) HandlePostNfticket(c *fiber.Ctx) error {
	var nfticket *types.Nfticket
	if err := c.BodyParser(&nfticket); err != nil {
		return errors.ErrBadRequest()
	}
	insertedNfticket, err := h.store.Nfticket.InsertNfticket(c.Context(), nfticket)
	if err != nil {
		return errors.ErrBadRequest()
	}
	return c.JSON(insertedNfticket)
}
