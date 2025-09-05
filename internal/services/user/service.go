package user

import (
	"context"

	"github.com/I-Van-Radkov/messenger/internal/models"
	"github.com/I-Van-Radkov/messenger/internal/repository"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (u *UserService) GetByID(ctx context.Context, id int64) (*models.User, error) {
	return u.userRepo.GetByID(ctx, id)
}

func (u *UserService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return u.userRepo.GetByEmail(ctx, email)
}

func (u *UserService) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return u.userRepo.GetByUsername(ctx, username)
}

func (u *UserService) Create(ctx context.Context, user *models.User) (int64, error) {
	return u.userRepo.Create(ctx, user)
}
