package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/greatdaveo/privycode-server/config"
	"github.com/greatdaveo/privycode-server/internal/middleware"
	"github.com/greatdaveo/privycode-server/internal/models"
)

func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUserFromContext(r)

	if user == nil {
		http.Error(w, "❌ Unauthorized", http.StatusUnauthorized)
		return
	}

	var links []models.ViewerLink

	if err := config.DB.Where("user_id = ?", user.ID).Find(&links).Error; err != nil {
		http.Error(w, "❌ Failed to fetch links", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(links)
}
