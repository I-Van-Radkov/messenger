package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

func hashPassword(password string, salt []byte) []byte {
	hashedPassword := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return hashedPassword
}

func generateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)

	return salt, err
}

func HashPasswordBase64(password string) (string, error) {
	salt, err := generateSalt()
	if err != nil {
		return "", err
	}

	hashedPassword := hashPassword(password, salt)

	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	hashedPasswordBase64 := base64.StdEncoding.EncodeToString(hashedPassword)

	return fmt.Sprintf("%s.%s", saltBase64, hashedPasswordBase64), nil
}

func VerifyPassword(newPasswordString, realHashedPasswordBase64 string) (bool, error) {
	parts := strings.Split(realHashedPasswordBase64, ".")

	salt, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, err
	}
	hashedRealPassword, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, err
	}

	hashedNewPassword := hashPassword(newPasswordString, salt)

	if len(hashedNewPassword) != len(hashedRealPassword) {
		return false, nil
	}
	if subtle.ConstantTimeCompare(hashedNewPassword, hashedRealPassword) == 0 {
		return false, nil
	}

	return true, nil
}
