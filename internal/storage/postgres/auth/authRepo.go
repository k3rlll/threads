package auth

import (
	"context"
	"main/pkg/customerrors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepo struct {
	pool *pgxpool.Pool
}

func NewAuthRepo(pool *pgxpool.Pool) *AuthRepo {
	return &AuthRepo{pool: pool}
}

func (r *AuthRepo) CreateUser(ctx context.Context, userID uuid.UUID, email, username, passwordHash string) (string, error) {
	tag, err := r.pool.Exec(ctx, "INSERT INTO users (id, email, username, password_hash) VALUES ($1, $2, $3, $4)",
		userID, email, username, passwordHash)
	if err != nil {
		return "", err
	}
	if tag.RowsAffected() != 1 {
		return "", customerrors.ErrNoTagsAffected
	}
	return userID.String(), nil
}

func (r *AuthRepo) GetUserByLogin(ctx context.Context, login string) (string, string, error) {
	
}

