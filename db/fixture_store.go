package db

// import (
// 	"context"
// 	"rockstake-core/types"

// 	"github.com/uptrace/bun"
// )

// type FixtureStore interface {
// 	InsertFixture(context.Context, *types.Fixture) (*types.Fixture, error)
// 	GetFixtureById(context.Context, int) *types.F
// }

// type SupabaseFixtureStore struct {
// 	client *bun.DB
// }

// func NewSupabaseFixtureStore(client *bun.DB) *SupabaseFixtureStore {
// 	return &SupabaseFixtureStore{
// 		client: client,
// 	}
// }

// func (s *SupabaseFixtureStore) InsertFixture(ctx context.Context, fixture *types.Fixture) (*types.Fixture, error) {
// 	_, err := s.client.NewInsert().Model(fixture).Exec(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return fixture, nil
// }
