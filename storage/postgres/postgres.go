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
	attribute  storage.AttributeI
	position   storage.PositionI
	posattr    storage.PosAttrI
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

func (s *Store) Attribute() storage.AttributeI {
	if s.attribute == nil {
		s.attribute = NewAttributeRepo(s.db)
	}

	return s.attribute
}

func (s *Store) Position() storage.PositionI {
	if s.position == nil {
		s.position = NewPositionRepo(s.db)
	}

	return s.position
}

func (s *Store) PosAttr() storage.PosAttrI {
	if s.posattr == nil {
		s.posattr = NewPosAttrRepo(s.db)
	}

	return s.posattr
}
