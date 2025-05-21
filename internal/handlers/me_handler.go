package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/greatdaveo/privycode-server/internal/middleware"
)

func MeHandler(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)
	if user == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"github_username": user.GitHubUsername,
		"email":           user.Email,
	})
}
