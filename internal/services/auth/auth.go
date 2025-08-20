package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/I-Van-Radkov/messenger/internal/config"
	"github.com/I-Van-Radkov/messenger/internal/models"
	"github.com/I-Van-Radkov/messenger/internal/repository"
	"github.com/I-Van-Radkov/messenger/internal/utils"
)

var (
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrUserNotFound          = errors.New("user not found")
	ErrInvalidPassword       = errors.New("invalid password")
	ErrInternalServer        = errors.New("internal server error")
)

type AuthService struct {
	userRepo repository.UserRepository
	cfg      *config.AuthConfig
}

func NewService(userRepo repository.UserRepository, cfg *config.AuthConfig) AuthService {
	return AuthService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s AuthService) Register(email, username, password string) (string, error) {
	ctx := context.Background()

	// Проверка на то, не существует ли пользователя по email
	_, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil {
		return "", ErrEmailAlreadyExists
	}
	if err != sql.ErrNoRows {
		return "", fmt.Errorf("%w: %v", ErrInternalServer, err)
	}

	// Проверка на уникальность username
	_, err = s.userRepo.GetByUsername(ctx, username)
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

	userId, err := s.userRepo.Create(ctx, user)
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

func (s AuthService) Login(email, password string) (string, error) {
	ctx := context.Background()
	var user *models.User

	// Найти поьзователя, если он действительно существует
	user, err := s.userRepo.GetByEmail(ctx, email)
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
