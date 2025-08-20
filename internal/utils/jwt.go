package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SignToken(userId int64, jwtSecret string, jwtExpiresIn time.Duration) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"user_id": userId,
		"iat":     now.Unix(),
		"nbf":     now.Unix(),
		"exp":     jwtExpiresIn,
	}

	tokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := tokenString.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
