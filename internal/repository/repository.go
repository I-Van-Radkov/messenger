package repository

import (
	"context"

	"github.com/I-Van-Radkov/messenger/internal/models"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (int64, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
}

type MessageRepository interface {
}
