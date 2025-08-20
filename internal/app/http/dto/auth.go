package dto

import "time"

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID         int       `json:"id"`
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	CreatedAt  time.Time `json:"created_at"`
	AvatarUrl  string    `json:"avatar_url"`
	LastOnline time.Time `json:"last_online"`
}
