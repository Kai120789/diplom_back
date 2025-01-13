package storage

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type Storage struct {
	UserStore *UserStorage
}

func Connect(DBUri string) (*pgxpool.Pool, error) {
	db, err := pgxpool.Connect(context.Background(), DBUri)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewPosgtresStorage(dbConn *pgxpool.Pool, log *zap.Logger) Storage {
	return Storage{
		UserStore: NewUserStorage(dbConn, log),
	}
}
