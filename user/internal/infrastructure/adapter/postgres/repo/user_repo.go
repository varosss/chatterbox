package repo

import (
	"chatterbox/user/internal/domain/entity"
	"chatterbox/user/internal/domain/valueobject"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserPgxRepo struct {
	db *pgxpool.Pool
}

func NewUserPgxRepo(db *pgxpool.Pool) *UserPgxRepo {
	return &UserPgxRepo{db: db}
}

func (r *UserPgxRepo) Save(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (
			id,
			email,
			username,
			display_name,
			password_hash,
			status
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (id) DO UPDATE SET
			email = EXCLUDED.email,
			username = EXCLUDED.username,
			display_name = EXCLUDED.display_name,
			password_hash = EXCLUDED.password_hash,
			status = EXCLUDED.status
	`

	_, err := r.db.Exec(
		ctx,
		query,
		user.ID().String(),
		user.Email().String(),
		user.Username(),
		user.DisplayName(),
		user.PasswordHash().String(),
		user.Status().Int(),
	)

	return err
}

func (r *UserPgxRepo) FindByID(
	ctx context.Context,
	id valueobject.UserID,
) (*entity.User, error) {
	query := `
		SELECT id, email, username, display_name, password_hash, status
		FROM users
		WHERE id = $1
		LIMIT 1
	`

	return r.scanOne(ctx, query, id.String())
}

func (r *UserPgxRepo) FindByEmail(
	ctx context.Context,
	email valueobject.Email,
) (*entity.User, error) {
	query := `
		SELECT id, email, username, display_name, password_hash, status
		FROM users
		WHERE email = $1
		LIMIT 1
	`

	return r.scanOne(ctx, query, email.String())
}

func (r *UserPgxRepo) ExistsByEmail(
	ctx context.Context,
	email valueobject.Email,
) bool {
	var exists bool

	err := r.db.QueryRow(
		ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`,
		email.String(),
	).Scan(&exists)

	return err == nil && exists
}

func (r *UserPgxRepo) List(
	ctx context.Context,
	userIDs []valueobject.UserID,
	limit int,
	offset int,
) ([]*entity.User, error) {
	ids := make([]string, len(userIDs))
	for i, id := range userIDs {
		ids[i] = id.String()
	}

	var args []any

	query := `
		SELECT id, email, username, display_name, password_hash, status
		FROM users
	`

	if len(ids) > 0 {
		query += " WHERE id = ANY($1) LIMIT ($2) OFFSET ($3)"
		args = append(args, ids, limit, offset)
	} else {
		query += " LIMIT ($1) OFFSET ($2)"
		args = append(args, limit, offset)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.User

	for rows.Next() {
		user, err := r.scanUser(rows)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, rows.Err()
}

func (r *UserPgxRepo) scanOne(
	ctx context.Context,
	query string,
	arg any,
) (*entity.User, error) {
	row, err := r.db.Query(ctx, query, arg)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	if !row.Next() {
		return nil, fmt.Errorf("user not found")
	}

	return r.scanUser(row)
}

func (r *UserPgxRepo) scanUser(rows pgx.Rows) (*entity.User, error) {
	var temp struct {
		ID           string
		Email        string
		Username     string
		DisplayName  string
		PasswordHash string
		Status       int
	}

	err := rows.Scan(
		&temp.ID,
		&temp.Email,
		&temp.Username,
		&temp.DisplayName,
		&temp.PasswordHash,
		&temp.Status,
	)
	if err != nil {
		return nil, err
	}

	id, err := valueobject.ParseUserID(temp.ID)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}

	email, err := valueobject.NewEmail(temp.Email)
	if err != nil {
		return nil, fmt.Errorf("parse email: %w", err)
	}

	passwordHash, err := valueobject.NewPasswordHash(temp.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("parse password hash: %w", err)
	}

	return entity.UserFromPrimitives(
		id,
		email,
		temp.Username,
		temp.DisplayName,
		passwordHash,
		valueobject.Status(temp.Status),
	), nil
}
