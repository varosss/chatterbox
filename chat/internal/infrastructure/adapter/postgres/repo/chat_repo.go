package repo

import (
	"chatterbox/chat/internal/domain/entity"
	"chatterbox/chat/internal/domain/valueobject"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatPgxRepo struct {
	db *pgxpool.Pool
}

func NewChatPgxRepo(db *pgxpool.Pool) *ChatPgxRepo {
	return &ChatPgxRepo{db: db}
}

func (r *ChatPgxRepo) Save(ctx context.Context, chat *entity.Chat) error {
	_, err := r.db.Exec(
		ctx,
		`INSERT INTO chats (
			id,
			participant_ids,
			display_name
		)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE SET
			participant_ids = EXCLUDED.participant_ids,
			display_name = EXCLUDED.display_name
		`,
		chat.ID().String(),
		chat.ParticipantIDsAsStrings(),
		chat.DisplayName(),
	)
	return err
}

func (r *ChatPgxRepo) FindByID(ctx context.Context, chatID valueobject.ChatID) (*entity.Chat, error) {
	query := `SELECT * FROM chats WHERE id = $1 LIMIT 1`
	var temp struct {
		ID             string
		ParticipantIDs []pgtype.UUID
		DisplayName    string
	}

	err := r.db.QueryRow(ctx, query, chatID.String()).Scan(&temp.ID, &temp.ParticipantIDs, &temp.DisplayName)
	if err != nil {
		return nil, fmt.Errorf("failed to find chat by chat ID: %w", err)
	}

	participantIDs := make([]valueobject.UserID, len(temp.ParticipantIDs))
	for i, elem := range temp.ParticipantIDs {
		parsed, err := valueobject.ParseUserID(elem.String())
		if err != nil {
			return nil, fmt.Errorf("failed to parse participant ID: %w", err)
		}
		participantIDs[i] = parsed
	}

	parsedChatID, err := valueobject.ParseChatID(temp.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse chat ID: %w", err)
	}

	return entity.ChatFromPrimitives(
		parsedChatID,
		participantIDs,
		temp.DisplayName,
	), nil
}

func (r *ChatPgxRepo) List(
	ctx context.Context,
	participantID valueobject.UserID,
) ([]*entity.Chat, error) {
	query := "SELECT * FROM chats WHERE $1 = ANY(participant_ids) LIMIT 500"
	rows, err := r.db.Query(ctx, query, participantID.String())
	if err != nil {
		return nil, err
	}

	var chats []*entity.Chat
	for rows.Next() {
		var temp struct {
			ID             string
			ParticipantIDs []pgtype.UUID
			DisplayName    string
		}

		rows.Scan(&temp.ID, &temp.ParticipantIDs, &temp.DisplayName)

		participantIDs := make([]valueobject.UserID, len(temp.ParticipantIDs))
		for i, elem := range temp.ParticipantIDs {
			parsed, err := valueobject.ParseUserID(elem.String())
			if err != nil {
				return nil, fmt.Errorf("failed to parse participant ID: %w", err)
			}
			participantIDs[i] = parsed
		}

		parsedChatID, err := valueobject.ParseChatID(temp.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to parse chat ID: %w", err)
		}

		chats = append(chats, entity.ChatFromPrimitives(
			parsedChatID,
			participantIDs,
			temp.DisplayName,
		))
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}
