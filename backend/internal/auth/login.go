package auth

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/giuszeppe/gatp-atc-2025/backend/internal/stores"
	"golang.org/x/crypto/bcrypt"
)

func Login(userStore stores.UserStore, username, password string) bool {
	// Retrieve the hashed password from the database
	user, err := userStore.GetUserWithUsername(username)
	if err != nil {
		return false
	}

	// Compare the provided password with the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func RandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
