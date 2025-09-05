package repository

import (
	"context"

	"github.com/I-Van-Radkov/messenger/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (int64, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
}

type MessageRepository interface {
	Create(ctx context.Context, message *models.Message) (int64, error)
	UpdateStatus(ctx context.Context, msgId int64, status string)

	GetMessagesByDialogID(ctx context.Context, dialogID int64, limit, offset int) ([]*models.Message, error)
	GetLastMessagesByDialogID(ctx context.Context, dialogIDs []int64) (map[int64]*models.Message, error)
}

type ChatRepository interface {
	GetChatsByUserID(ctx context.Context, userID int64, limit, offset int) ([]*models.Chat, error)
}
