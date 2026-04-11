package repo

import (
	"chatterbox/chat/internal/domain/entity"
	"chatterbox/chat/internal/domain/valueobject"
	"context"
	"fmt"
	"time"

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
		`INSERT INTO messages (
			id,
			chat_id,
			sender_id,
			text,
			created_at
		)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE SET
			text = EXCLUDED.text
		`,
		m.ID().String(),
		m.ChatID().String(),
		m.SenderID().String(),
		m.Text(),
		m.CreatedAt(),
	)
	return err
}

func (r *MessagePgxRepo) List(
	ctx context.Context,
	chatID valueobject.ChatID,
) ([]*entity.Message, error) {
	query := "SELECT * FROM messages WHERE chat_id = $1 LIMIT 500"
	rows, err := r.db.Query(ctx, query, chatID.String())
	if err != nil {
		return nil, err
	}

	var messages []*entity.Message
	for rows.Next() {
		var temp struct {
			ID        string
			ChatID    string
			SenderID  string
			Text      string
			CreatedAt time.Time
		}

		rows.Scan(
			&temp.ID,
			&temp.ChatID,
			&temp.SenderID,
			&temp.Text,
			&temp.CreatedAt,
		)

		parsedMessageID, err := valueobject.ParseMessageID(temp.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse message ID: %w", err)
		}

		parsedChatID, err := valueobject.ParseChatID(temp.ChatID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse chat ID: %w", err)
		}

		parsedSenderID, err := valueobject.ParseUserID(temp.SenderID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse sender ID: %w", err)
		}

		messages = append(messages, entity.MessageFromPrimitives(
			parsedMessageID,
			parsedSenderID,
			parsedChatID,
			temp.Text,
			temp.CreatedAt,
		))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
