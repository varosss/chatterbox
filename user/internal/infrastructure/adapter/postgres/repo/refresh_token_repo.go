package repo

import (
	"chatterbox/user/internal/domain/entity"
	"chatterbox/user/internal/domain/valueobject"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenPgxRepo struct {
	db *pgxpool.Pool
}

func NewRefreshTokenPgxRepo(db *pgxpool.Pool) *RefreshTokenPgxRepo {
	return &RefreshTokenPgxRepo{db: db}
}

func (r *RefreshTokenPgxRepo) Save(
	ctx context.Context,
	token *entity.RefreshToken,
) error {
	query := `
		INSERT INTO refresh_tokens (
			id,
			user_id,
			expires_at,
			revoked
		)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			revoked = EXCLUDED.revoked,
			expires_at = EXCLUDED.expires_at
	`

	_, err := r.db.Exec(
		ctx,
		query,
		token.ID().String(),
		token.UserID().String(),
		token.ExpiresAt(),
		token.IsRevoked(),
	)

	return err
}

func (r *RefreshTokenPgxRepo) Get(
	ctx context.Context,
	id valueobject.TokenID,
) (*entity.RefreshToken, error) {
	query := `
		SELECT id, user_id, expires_at, revoked
		FROM refresh_tokens
		WHERE id = $1
		LIMIT 1
	`

	var temp struct {
		ID        string
		UserID    string
		ExpiresAt time.Time
		Revoked   bool
	}

	err := r.db.QueryRow(ctx, query, id.String()).Scan(
		&temp.ID,
		&temp.UserID,
		&temp.ExpiresAt,
		&temp.Revoked,
	)
	if err != nil {
		return nil, err
	}

	tokenID, err := valueobject.ParseTokenID(temp.ID)
	if err != nil {
		return nil, fmt.Errorf("parse token id: %w", err)
	}

	userID, err := valueobject.ParseUserID(temp.UserID)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}

	return entity.RefreshTokenFromPrimitives(
		tokenID,
		userID,
		temp.ExpiresAt,
		temp.Revoked,
	), nil
}

func (r *RefreshTokenPgxRepo) Revoke(
	ctx context.Context,
	id valueobject.TokenID,
) error {
	_, err := r.db.Exec(
		ctx,
		`UPDATE refresh_tokens SET revoked = TRUE WHERE id = $1`,
		id.String(),
	)

	return err
}
