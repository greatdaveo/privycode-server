package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/greatdaveo/privycode-server/config"
	"github.com/greatdaveo/privycode-server/internal/models"
)

type contextKey string

const userCtxKey = contextKey("user")

// To extract user from Authorization header
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if !strings.HasPrefix(authHeader, "Bearer") {
			http.Error(w, "❌ Unauthorized: Missed token", http.StatusUnauthorized)
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		var user models.User
		if err := config.DB.Where("git_hub_token = ?", token).First(&user).Error; err != nil {
			http.Error(w, "❌ Unauthorized: Invalid token", http.StatusUnauthorized)
		}

		ctx := context.WithValue(r.Context(), userCtxKey, &user)
		next(w, r.WithContext(ctx))
	}
}

// To extract the user from the context
func GetUserFromContext(r *http.Request) *models.User {
	user, ok := r.Context().Value(userCtxKey).(*models.User)
	if !ok {
		return nil
	}

	return user
}
