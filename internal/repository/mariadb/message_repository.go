package mariadb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/I-Van-Radkov/messenger/internal/models"
)

type MessageRepo struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepo {
	return &MessageRepo{
		db: db,
	}
}

type messageDB struct {
	ID           sql.NullInt64
	DialogID     sql.NullInt64
	SenderID     sql.NullInt64
	RecipientID  sql.NullInt64
	Content      sql.NullString
	CreatedAt    sql.NullTime
	IsReplyToMsg sql.NullBool
	ReplyToMsgID sql.NullInt64
	Status       sql.NullString
}

func (r *MessageRepo) toDomain(dbMsg *messageDB) *models.Message {
	msg := &models.Message{
		ID:           dbMsg.ID.Int64,
		DialogID:     dbMsg.DialogID.Int64,
		SenderID:     dbMsg.SenderID.Int64,
		RecipientID:  dbMsg.RecipientID.Int64,
		Content:      dbMsg.Content.String,
		CreatedAt:    dbMsg.CreatedAt.Time,
		IsReplyToMsg: dbMsg.IsReplyToMsg.Bool,
		Status:       dbMsg.Status.String,
	}

	if dbMsg.ReplyToMsgID.Valid {
		msg.ReplyToMsgID = dbMsg.ReplyToMsgID.Int64
	}

	return msg
}

func (r *MessageRepo) toDB(domainMsg *models.Message) *messageDB {
	return &messageDB{
		ID:           sql.NullInt64{Int64: domainMsg.ID, Valid: true},
		DialogID:     sql.NullInt64{Int64: domainMsg.DialogID, Valid: true},
		SenderID:     sql.NullInt64{Int64: domainMsg.SenderID, Valid: true},
		RecipientID:  sql.NullInt64{Int64: domainMsg.RecipientID, Valid: true},
		Content:      sql.NullString{String: domainMsg.Content, Valid: true},
		CreatedAt:    sql.NullTime{Time: domainMsg.CreatedAt, Valid: true},
		IsReplyToMsg: sql.NullBool{Bool: domainMsg.IsReplyToMsg, Valid: true},
		ReplyToMsgID: sql.NullInt64{Int64: domainMsg.ReplyToMsgID, Valid: true},
		Status:       sql.NullString{String: domainMsg.Status, Valid: true},
	}
}

func (r *MessageRepo) Create(ctx context.Context, message *models.Message) (int64, error) {
	query := `INSERT INTO messages (
				dialog_id,
				sender_id,
				recipient_id,
				content,
				created_at,
				is_reply_to_msg,
				reply_to_msg_id,
				status
	) VALUES (?,?,?,?,?,?,?,?)`

	dbMsg := r.toDB(message)

	result, err := r.db.ExecContext(ctx, query,
		dbMsg.DialogID,
		dbMsg.SenderID,
		dbMsg.RecipientID,
		dbMsg.Content,
		dbMsg.CreatedAt,
		dbMsg.IsReplyToMsg,
		dbMsg.ReplyToMsgID,
		dbMsg.Status,
	)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	return id, nil
}

func (r *MessageRepo) UpdateStatus(ctx context.Context, msgId int64, status string) {
	query := "UPDATE messages SET status = ? WHERE id = ?"

	r.db.ExecContext(ctx, query, status, msgId)
}

func (r *MessageRepo) GetMessagesByDialogID(ctx context.Context, dialogID int64, limit, offset int) ([]*models.Message, error) {
	query := `SELECT
				id,
				dialog_id,
				sender_id,
				recipient_id,
				content,
				created_at,
				is_reply_to_msg,
				reply_to_msg_id,
				status
			  FROM messages
			  WHERE dialog_id = ?
			  ORDER BY created_at DESC, id DESC
			  LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, dialogID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query messages: %w", err)
	}
	defer rows.Close()

	var messages []*models.Message
	for rows.Next() {
		var dbMsg messageDB
		err := rows.Scan(
			&dbMsg.ID,
			&dbMsg.DialogID,
			&dbMsg.SenderID,
			&dbMsg.RecipientID,
			&dbMsg.Content,
			&dbMsg.CreatedAt,
			&dbMsg.IsReplyToMsg,
			&dbMsg.ReplyToMsgID,
			&dbMsg.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

		messages = append(messages, r.toDomain(&dbMsg))
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return messages, nil
}

func (r *MessageRepo) GetLastMessagesByDialogID(ctx context.Context, dialogIDs []int64) (map[int64]*models.Message, error) {

	return nil, nil
}
