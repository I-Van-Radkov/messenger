package auth

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/I-Van-Radkov/messenger/internal/config"
	"github.com/I-Van-Radkov/messenger/internal/models"
	"github.com/I-Van-Radkov/messenger/internal/utils"
)

type UserProvider interface {
	Create(ctx context.Context, user *models.User) (int64, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
}

type AuthService struct {
	userService UserProvider
	cfg         *config.AuthConfig
}

func NewService(userService UserProvider, cfg *config.AuthConfig) AuthService {
	return AuthService{
		userService: userService,
		cfg:         cfg,
	}
}

func (s AuthService) Register(ctx context.Context, email, username, password string) (string, error) {
	// Проверка на то, не существует ли пользователя по email
	_, err := s.userService.GetByEmail(ctx, email)
	if err == nil {
		return "", ErrEmailAlreadyExists
	}
	if err != sql.ErrNoRows {
		return "", fmt.Errorf("%w: %v", ErrInternalServer, err)
	}

	// Проверка на уникальность username
	_, err = s.userService.GetByUsername(ctx, username)
	if err == nil {
		return "", ErrUsernameAlreadyExists
	}
	if err != sql.ErrNoRows {
		return "", fmt.Errorf("%w: %v", ErrInternalServer, err)
	}

	// Хеширование пароля
	hashedPasswordBase64, err := utils.HashPasswordBase64(password)
	if err != nil {
		return "", fmt.Errorf("%w: failed to hash password", ErrInternalServer)
	}

	// Сохранение данных в БД
	user := &models.User{
		Email:        email,
		Username:     username,
		PasswordHash: hashedPasswordBase64,
		CreatedAt:    time.Now(),
	}

	userId, err := s.userService.Create(ctx, user)
	if err != nil {
		return "", fmt.Errorf("%w: failed to create user", ErrInternalServer)
	}

	// Генерация JWT токена
	tokenString, err := utils.SignToken(userId, s.cfg.JwtSecret, s.cfg.JwtExpiresIn)
	if err != nil {
		return "", fmt.Errorf("%w: failed to generate token", ErrInternalServer)
	}

	// Возврат JWT токена
	return tokenString, nil
}

func (s AuthService) Login(ctx context.Context, email, password string) (string, error) {
	var user *models.User

	// Найти поьзователя, если он действительно существует
	user, err := s.userService.GetByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrUserNotFound
		}
		return "", fmt.Errorf("%w: failed to get user", ErrInternalServer)
	}

	// Подтвердить пароль
	ok, err := utils.VerifyPassword(password, user.PasswordHash)
	if err != nil {
		return "", fmt.Errorf("%w: failed to verify password", ErrInternalServer)
	}
	if !ok {
		return "", ErrInvalidPassword
	}

	// Генерация JWT токена
	tokenString, err := utils.SignToken(user.ID, s.cfg.JwtSecret, s.cfg.JwtExpiresIn)
	if err != nil {
		return "", fmt.Errorf("%w: failed to generate token", ErrInternalServer)
	}

	// Возврат JWT токена
	return tokenString, nil
}
