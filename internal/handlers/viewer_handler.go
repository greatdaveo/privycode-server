package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

func ViewerAccessHandler(w http.ResponseWriter, r *http.Request) {
	// To extract the token from the URL
	token := strings.TrimPrefix(r.URL.Path, "/view/")
	if token == "" {
		http.Error(w, "Missing token: ", http.StatusBadRequest)
		return
	}

	dbInstance := config.DB
	var link models.ViewerLink

	result := dbInstance.Where("token = ?", token).First(&link)
	if result.Error != nil {
		http.Error(w, "❌ Invalid or expired link", http.StatusNotFound)
		return
	}

	// To check if the link has expired
	if time.Now().After(link.ExpiresAt) {
		http.Error(w, "❌ Link has expired", http.StatusForbidden)
		return
	}

	// To check the max views (Optional)
	if link.MaxViews > 0 && link.ViewCount >= link.MaxViews {
		http.Error(w, "❌ View limit reached", http.StatusForbidden)
		return
	}

	// To increase the view count
	link.ViewCount++
	dbInstance.Save(&link)

	fmt.Fprintf(w, "✅ Access granted to repo: %s", link.RepoName)
}
