package repo

import (
	"chatterbox/chat/internal/domain/entity"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MessagePgxRepo struct {
	db *pgxpool.Pool
}

func NewMessagePgxRepo(db *pgxpool.Pool) *MessagePgxRepo {
	return &MessagePgxRepo{db: db}
}

func (r *MessagePgxRepo) Save(ctx context.Context, m *entity.Message) error {
	_, err := r.db.Exec(
		ctx,
		`INSERT INTO messages (id, chat_id, sender_id, text, created_at) VALUES ($1, $2, $3, $4, $5)`,
		m.ID().String(),
		m.ChatID().String(),
		m.SenderID().String(),
		m.Text(),
		m.CreatedAt(),
	)
	return err
}
