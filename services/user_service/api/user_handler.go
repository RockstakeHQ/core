package api

import (
	"betcube-engine/types"

	datastore "betcube-engine/datastore"
	errors "betcube-engine/errors"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	store *datastore.SupabaseStore
}

func NewUserHandler(s *datastore.SupabaseStore) *UserHandler {
	return &UserHandler{
		store: s,
	}
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var user *types.User
	if err := c.BodyParser(&user); err != nil {
		return err
	}
	insertedUser, err := h.store.User.InsertUser(c.Context(), user)
	if err != nil {
		return errors.ErrBadRequest()
	}
	return c.JSON(insertedUser)
}

func (h *UserHandler) HandleGetUserByWalletAddress(c *fiber.Ctx) error {
	walletAddress := c.Params("wallet_address")
	user, err := h.store.User.GetUserByWalletAddress(c.Context(), walletAddress)
	if err != nil {
		return errors.ErrNotResourceNotFound("wallet_address")
	}
	return c.JSON(user)
}

func (h *UserHandler) HandleGetUserByUserId(c *fiber.Ctx) error {
	client_id := c.Params("user_id")
	user, err := h.store.User.GetUserByUserId(c.Context(), client_id)
	if err != nil {
		return errors.ErrNotResourceNotFound("user_id")
	}
	return c.JSON(user)
}
