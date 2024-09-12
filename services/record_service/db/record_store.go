package db

import (
	"betcube_engine/types"
	"context"
	"log"

	"github.com/uptrace/bun"
)

type RecordStore interface {
	InsertRecord(context.Context, *types.Record) (*types.Record, error)
	GetRecordsByUser(context.Context, string) (*[]types.Record, error)
}

type SupabaseRecordStore struct {
	client *bun.DB
}

func NewSupabaseRecordStore(client *bun.DB) *SupabaseRecordStore {
	return &SupabaseRecordStore{
		client: client,
	}
}

func (s *SupabaseRecordStore) InsertRecord(ctx context.Context, record *types.Record) (*types.Record, error) {
	_, err := s.client.NewInsert().Model(record).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (s *SupabaseRecordStore) GetRecordsByUser(ctx context.Context, user string) (*[]types.Record, error) {
	var records []types.Record
	query := s.client.NewSelect().Model(&records).Where("wallet_address = ?", user)
	err := query.Scan(ctx, &records)
	if err != nil {
		return nil, err
	}
	log.Printf("Found %d records for user %s", len(records), user)
	return &records, nil
}
