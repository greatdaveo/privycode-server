package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/greatdaveo/privycode-server/config"
	"github.com/greatdaveo/privycode-server/internal/middleware"
	"github.com/greatdaveo/privycode-server/internal/models"
	"github.com/greatdaveo/privycode-server/internal/utils"
)

type ViewerLinkRequest struct {
	RepoName  string `json:"repo_name"`
	ExpiresIn int    `json:"expires_in_days"`
	MaxViews  int    `json:"max_views"`
}

func GenerateViewerLinkHandler(w http.ResponseWriter, r *http.Request) {
	// Temp User ID
	user := middleware.GetUserFromContext(r)
	// fmt.Println("User: ", user)

	if user == nil {
		http.Error(w, "❌ Unauthorized", http.StatusUnauthorized)
		return
	}

	var req ViewerLinkRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.RepoName == "" {
		http.Error(w, "❌ Invalid input", http.StatusBadRequest)
		return
	}

	token := utils.GenerateToken()
	// To calculate expiring date
	days := req.ExpiresIn
	if days <= 0 {
		days = 3 // to set default expiration to 3 days
	}
	expiration := time.Now().Add(time.Duration(days) * 24 * time.Hour)

	link := models.ViewerLink{
		RepoName:  req.RepoName,
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiration,
		MaxViews:  req.MaxViews,
		ViewCount: 0,
	}

	// To ensure the repository exist before saving
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s", user.GitHubUsername, req.RepoName)
	reqGitHub, _ := http.NewRequest("GET", apiURL, nil)
	reqGitHub.Header.Set("Authorization", "token "+user.GitHubToken)
	reqGitHub.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(reqGitHub)

	if err != nil || resp.StatusCode != 200 {
		http.Error(w, "❌ Repository not found or inaccessible", http.StatusNotFound)
		return
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

	// To get the owner of the repo
	var user models.User
	err := dbInstance.First(&user, link.UserID).Error
	if err != nil {
		http.Error(w, "❌ User not found", http.StatusInternalServerError)
		return
	}

	// To build GitHub API request
	client := &http.Client{}
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents", user.GitHubUsername, link.RepoName)

	// fmt.Println("GitHub API URL:", apiURL)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		http.Error(w, "❌ Failed to build request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", "token "+user.GitHubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	// To send the request
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "❌ HTTP error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	// To read error response body
	if resp.StatusCode != 200 {
		errorBody, _ := io.ReadAll(resp.Body)
		http.Error(w, string(errorBody), resp.StatusCode)
		return
	}
	// To parse GitHub response
	var contents []struct {
		Name string `json:"name"`
		Type string `json:"type"`
		Path string `json:"path"`
		URL  string `json:"url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&contents); err != nil {
		http.Error(w, "❌ Failed to parse GitHub response", http.StatusInternalServerError)
		return
	}

	// To return the content list as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contents)

	// fmt.Fprintf(w, "✅ Access granted to repo: %s", link.RepoName)
}

func ViewFileHandler(w http.ResponseWriter, r *http.Request) {
	segments := strings.Split(strings.TrimPrefix(r.URL.Path, "/view-files/"), "/")

	if len(segments) < 1 {
		http.Error(w, "Invalid viewer URL", http.StatusBadRequest)
		return
	}

	token := segments[0]

	// fmt.Println("Extracted Token:", token)

	path := r.URL.Query().Get("path")

	if token == "" || path == "" {
		http.Error(w, "❌ Missing token or path", http.StatusBadRequest)
		return
	}

	dbInstance := config.DB

	// To get the viewer link
	var link models.ViewerLink
	if err := dbInstance.Where("token = ?", token).First(&link).Error; err != nil {
		http.Error(w, "❌ Invalid link", http.StatusNotFound)
		return
	}

	// To check expiration
	if time.Now().After(link.ExpiresAt) {
		http.Error(w, "❌ Link expired", http.StatusForbidden)
		return
	}
	// To check view limits
	if link.MaxViews > 0 && link.ViewCount >= link.MaxViews {
		http.Error(w, "❌ View limit reached", http.StatusForbidden)
		return
	}

	var user models.User
	if err := dbInstance.First(&user, link.UserID).Error; err != nil {
		http.Error(w, "❌ User not found", http.StatusInternalServerError)
		return
	}

	// To request file content from GitHub
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", user.GitHubUsername, link.RepoName, path)
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("Authorization", "token "+user.GitHubToken)
	req.Header.Set("Accept", "application/vnd.github.v3.raw")

	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil || response.StatusCode != 200 {
		body, _ := io.ReadAll(response.Body)
		http.Error(w, fmt.Sprintf("❌ GitHub error: %s", body), response.StatusCode)
		return
	}

	defer response.Body.Close()

	content, _ := io.ReadAll(response.Body)
	w.Header().Set("Content-Type", "text/plain")
	w.Write(content)
}

func ViewerFolderHandler(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.URL.Path, "/view-folder/")
	if token == "" {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	// fmt.Println("Folder View Token:", token)

	path := r.URL.Query().Get("path")

	dbInstance := config.DB

	// To look up viewer link
	var link models.ViewerLink
	if err := dbInstance.Where("token = ?", token).First(&link).Error; err != nil {
		http.Error(w, "Invalid viewer link", http.StatusNotFound)
		return
	}

	// To check expiration and view count
	if time.Now().After(link.ExpiresAt) || (link.MaxViews > 0 && link.ViewCount >= link.MaxViews) {
		http.Error(w, "Link expired or view limit reached", http.StatusForbidden)
		return
	}

	// To get the user
	var user models.User
	if err := dbInstance.First(&user, link.UserID).Error; err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// To build the GitHub API folder URL
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", user.GitHubUsername, link.RepoName, path)
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("Authorization", "token "+user.GitHubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		http.Error(w, fmt.Sprintf("GitHub folder fetch error: %s", body), resp.StatusCode)
		return
	}
	defer resp.Body.Close()

	// To return the folder content
	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}

func ViewUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.URL.Path, "/view-info/")
	if token == "" {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	var link models.ViewerLink
	if err := config.DB.Where("token = ?", token).First(&link).Error; err != nil {
		http.Error(w, "Invalid or expired link", http.StatusNotFound)
		return
	}

	var user models.User
	if err := config.DB.First(&user, link.UserID).Error; err != nil {
		http.Error(w, "User not found", http.StatusInternalServerError)
		return
	}

	// To return the user GitHub username and repo name
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"github_username": user.GitHubUsername,
		"repo_name":       link.RepoName,
	})
}

func UpdateViewerLinkHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/update-link/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "❌ Invalid ID", http.StatusBadRequest)
		return
	}

	var link models.ViewerLink

	db := config.DB
	if err := db.First(&link, id).Error; err != nil {
		http.Error(w, "❌ Link not found", http.StatusNotFound)
		return
	}

	var payload struct {
		ExpiresInDays int `json:"expires_in_days"`
		MaxViews      int `json:"max_views"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "❌ Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if payload.ExpiresInDays > 0 {
		link.ExpiresAt = time.Now().Add(time.Duration(payload.ExpiresInDays) * 24 * time.Hour)
	}

	if payload.MaxViews > 0 {
		link.MaxViews = payload.MaxViews
	}

	if err := db.Save(&link).Error; err != nil {
		http.Error(w, "❌ Could not update link", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Viewer link updated successfully",
	})
}

func DeleteViewerLinkHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/delete-link/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, "❌ Invalid ID", http.StatusBadRequest)
		return
	}

	var link models.ViewerLink
	db := config.DB
	if err := db.First(&link, id).Error; err != nil {
		http.Error(w, "❌ Link not found", http.StatusNotFound)
		return
	}

	if err := db.Delete(&link).Error; err != nil {
		http.Error(w, "❌  Could not delete link", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Viewer link deleted",
	})
}
