package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/greatdaveo/privycode-server/config"
	"github.com/greatdaveo/privycode-server/internal/models"
	"github.com/greatdaveo/privycode-server/internal/utils"
)

type ViewerLinkRequest struct {
	RepoName  string `json:"repo_name"`
	ExpiresIn int    `json:"expires_in_days`
	MaxViews  int    `json:"max_views"`
}

func GenerateViewerLinkHandler(w http.ResponseWriter, r *http.Request) {
	// Temp User ID
	userID := uint(1)

	var req ViewerLinkRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.RepoName == "" {
		http.Error(w, "❌ Invalid input", http.StatusBadRequest)
		return
	}

	token := utils.GenerateToken()
	expiration := time.Now().Add(time.Duration(req.ExpiresIn) * 24 * time.Hour)

	link := models.ViewerLink{
		RepoName:  req.RepoName,
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiration,
		MaxViews:  req.MaxViews,
		ViewCount: 0,
	}

	if err := config.DB.Create(&link).Error; err != nil {
		http.Error(w, "❌ Could not create viewer link", http.StatusInternalServerError)
	}

	viewerURL := fmt.Sprintf("http://localhost:8080/view/%s", token)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"viewer_url": viewerURL,
	})
}


