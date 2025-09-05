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

type userDB struct {
	ID           sql.NullInt64
	Email        sql.NullString
	Username     sql.NullString
	PasswordHash sql.NullString
	CreatedAt    sql.NullTime
}

func (r *UserRepo) toDomain(dbUser *userDB) *models.User {
	return &models.User{
		ID:           dbUser.ID.Int64,
		Email:        dbUser.Email.String,
		Username:     dbUser.Username.String,
		PasswordHash: dbUser.PasswordHash.String,
		CreatedAt:    dbUser.CreatedAt.Time,
	}
}

func (r *UserRepo) toDB(domainUser *models.User) *userDB {
	return &userDB{
		ID:           sql.NullInt64{Int64: domainUser.ID, Valid: true},
		Email:        sql.NullString{String: domainUser.Email, Valid: true},
		Username:     sql.NullString{String: domainUser.Username, Valid: true},
		PasswordHash: sql.NullString{String: domainUser.PasswordHash, Valid: true},
		CreatedAt:    sql.NullTime{Time: domainUser.CreatedAt, Valid: true},
	}
}

func (r *UserRepo) Create(ctx context.Context, user *models.User) (int64, error) {
	query := "INSERT INTO users (email, username, password_hash, created_at) VALUES (?,?,?,?)"

	dbUser := r.toDB(user)

	result, err := r.db.ExecContext(ctx, query, dbUser.Email, dbUser.Username, dbUser.PasswordHash, dbUser.CreatedAt)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := "SELECT id, email, username, password_hash, created_at FROM users WHERE email = ?"

	var dbUser userDB

	err := r.db.QueryRowContext(ctx, query, email).Scan(&dbUser.ID, &dbUser.Email, &dbUser.Username, &dbUser.PasswordHash, &dbUser.CreatedAt)
	if err != nil {
		return nil, err
	}

	return r.toDomain(&dbUser), nil
}

func (ur *UserRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := "SELECT id, email, username, password_hash, created_at FROM users WHERE username = ?"

	var dbUser userDB

	err := ur.db.QueryRowContext(ctx, query, username).Scan(&dbUser.ID, &dbUser.Email, &dbUser.Username, &dbUser.PasswordHash, &dbUser.CreatedAt)
	if err != nil {
		return nil, err
	}

	return ur.toDomain(&dbUser), nil
}

func (ur *UserRepo) GetByID(ctx context.Context, id int64) (*models.User, error) {
	query := "SELECT id, email, username, password_hash, created_at FROM users WHERE id = ?"

	var dbUser userDB

	err := ur.db.QueryRowContext(ctx, query, id).Scan(&dbUser.ID, &dbUser.Email, &dbUser.Username, &dbUser.PasswordHash, &dbUser.CreatedAt)
	if err != nil {
		return nil, err
	}

	return ur.toDomain(&dbUser), nil
}
