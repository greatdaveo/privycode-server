package middleware

import (
	"context"
	"net/http"

	"github.com/greatdaveo/privycode-server/config"
	"github.com/greatdaveo/privycode-server/internal/models"
)

type contextKey string

const userCtxKey = contextKey("user")

// To extract user from Authorization header
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("github_token")
		if err != nil || cookie.Value == "" {
			http.Error(w, "❌ Unauthorized: Missing token", http.StatusUnauthorized)
			return
		}

		token := cookie.Value
		
		var user models.User
		// fmt.Println("user", user)

		if err := config.DB.Where("git_hub_token = ?", token).First(&user).Error; err != nil {
			// log.Println("❌ Auth failed:", err)
			http.Error(w, "❌ Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userCtxKey, &user)
		next.ServeHTTP(w, r.WithContext(ctx))
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
