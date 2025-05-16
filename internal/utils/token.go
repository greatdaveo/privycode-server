package utils

import "github.com/google/uuid"

// To create a new UUID token string
func GenerateToken() string {
	return uuid.New().String()
}
