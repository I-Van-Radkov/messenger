package dto

import "time"

type UserSearchResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type UserProfileResponse struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}
