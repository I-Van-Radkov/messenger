package mariadb

import (
	"context"
	"database/sql"

	"github.com/I-Van-Radkov/messenger/internal/models"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (ur *UserRepo) Create(ctx context.Context, user *models.User) (int64, error) {
	query := "INSERT INTO users (email, username, password_hash, created_at) VALUES (?,?,?,?)"
	result, err := ur.db.ExecContext(ctx, query, user.Email, user.Username, user.PasswordHash, user.CreatedAt)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (ur *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := "SELECT id, email, username, password_hash, created_at FROM users WHERE email = ?"

	err := ur.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	query := "SELECT id, email, username, password_hash, created_at FROM users WHERE username = ?"

	err := ur.db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
