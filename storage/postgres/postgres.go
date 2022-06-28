package postgres

import (
	"context"
	"position_service/config"
	"position_service/storage"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	db         *pgxpool.Pool
	profession storage.ProfessionI
}

func NewPostgres(psqlConnString string, cfg config.Config) (storage.StorageI, error) {
	// First set up the pgx connection pool
	config, err := pgxpool.ParseConfig(psqlConnString)
	if err != nil {
		return nil, err
	}

	config.AfterConnect = nil
	config.MaxConns = int32(cfg.PostgresMaxConnections)

	pool, err := pgxpool.ConnectConfig(context.Background(), config)

	return &Store{
		db: pool,
	}, err
}

func (s *Store) Profession() storage.ProfessionI {
	if s.profession == nil {
		s.profession = NewProfessionRepo(s.db)
	}

	return s.profession
}
