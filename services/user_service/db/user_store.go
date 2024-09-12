package db

import (
	"betcube-engine/types"
	"context"

	"github.com/uptrace/bun"
)

type UserStore interface {
	InsertUser(context.Context, *types.User) (*types.User, error)
	GetUserByWalletAddress(context.Context, string) (*types.User, error)
	GetUserByUserId(context.Context, string) (*types.User, error)
}

type SupabaseUserStore struct {
	client *bun.DB
}

func NewSupabaseUserStore(client *bun.DB) *SupabaseUserStore {
	return &SupabaseUserStore{
		client: client,
	}
}

func (s *SupabaseUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {
	_, err := s.client.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *SupabaseUserStore) GetUserByWalletAddress(ctx context.Context, walletAddress string) (*types.User, error) {
	var user types.User
	err := s.client.NewSelect().Model(&user).Where("wallet_address = ?", walletAddress).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *SupabaseUserStore) GetUserByUserId(ctx context.Context, userId string) (*types.User, error) {
	var user types.User
	err := s.client.NewSelect().Model(&user).Where("user_id = ?", userId).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
