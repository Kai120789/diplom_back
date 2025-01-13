package storage

import (
	"context"
	"materials/internal/apperrors"
	"materials/internal/dto"
	"materials/internal/models"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type UserStorage struct {
	log *zap.Logger
	db  *pgxpool.Pool
}

func NewUserStorage(db *pgxpool.Pool, log *zap.Logger) *UserStorage {
	return &UserStorage{db: db, log: log}
}

func (u *UserStorage) CreateUser(dto dto.RegistrationUser) error {
	sql := `INSERT INTO users (username, password) VALUES ($1, $2)`
	_, err := u.db.Exec(context.Background(), sql, dto.Username, dto.Password)
	if err != nil {
		u.log.Sugar().Errorf("create user", err)
		return err
	}
	return nil
}

func (u *UserStorage) GetUserByName(username string) (*models.User, error) {
	sql := `
		SELECT 
			id, 
			username, 
			password,
			created_at,
			updated_at
		FROM users WHERE username = $1`
	row := u.db.QueryRow(context.Background(), sql, username)
	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}
