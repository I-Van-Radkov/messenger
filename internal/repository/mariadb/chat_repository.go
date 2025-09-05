package mariadb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/I-Van-Radkov/messenger/internal/models"
)

type ChatRepo struct {
	db *sql.DB
}

func NewChatRepo(db *sql.DB) *ChatRepo {
	return &ChatRepo{
		db: db,
	}
}

type chatDB struct {
	ID            sql.NullInt64
	User1ID       sql.NullInt64
	User2ID       sql.NullInt64
	CreatedAt     sql.NullTime
	LastMessageID sql.NullInt64
}

func (r *ChatRepo) toDomain(dbChat *chatDB) *models.Chat {
	return &models.Chat{
		ID:            dbChat.ID.Int64,
		User1ID:       dbChat.User1ID.Int64,
		User2ID:       dbChat.User2ID.Int64,
		CreatedAt:     dbChat.CreatedAt.Time,
		LastMessageID: dbChat.LastMessageID.Int64,
	}
}

func (r *ChatRepo) toDB(domainChat *models.Chat) *chatDB {
	return &chatDB{
		ID:            sql.NullInt64{Int64: domainChat.ID, Valid: true},
		User1ID:       sql.NullInt64{Int64: domainChat.User1ID, Valid: true},
		User2ID:       sql.NullInt64{Int64: domainChat.User2ID, Valid: true},
		CreatedAt:     sql.NullTime{Time: domainChat.CreatedAt, Valid: true},
		LastMessageID: sql.NullInt64{Int64: domainChat.LastMessageID, Valid: true},
	}
}

func (r *ChatRepo) GetChatsByUserID(ctx context.Context, userID int64, limit, offset int) ([]*models.Chat, error) {
	query := `SELECT
				id,
				user1_id,
				user2_id,
				created_at,
				last_message_id
			FROM chats
			WHERE user1_id = ? OR user2_id = ?
			ORDER BY created_at DESC
			LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, userID, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query chats: %w", err)
	}
	defer rows.Close()

	var chats []*models.Chat
	for rows.Next() {
		var dbChat chatDB
		err := rows.Scan(
			&dbChat.ID,
			&dbChat.User1ID,
			&dbChat.User2ID,
			&dbChat.CreatedAt,
			&dbChat.LastMessageID,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan chat: %w", err)
		}

		chats = append(chats, r.toDomain(&dbChat))
	}

	return chats, nil
}
