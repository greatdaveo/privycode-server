package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/greatdaveo/privycode-server/config"
	"github.com/greatdaveo/privycode-server/internal/github"
	"github.com/greatdaveo/privycode-server/internal/models"
	"gorm.io/gorm"
)

func GitHubLoginHandler(w http.ResponseWriter, r *http.Request) {
	state := uuid.New().String()
	url := github.GetAuthURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

func GitHubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing code in URL", http.StatusBadRequest)
		return
	}

	token, err := github.ExchangeCodeForToken(code)
	if err != nil {
		http.Error(w, "❌ Failed to exchange code: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// To get the user info with the token
	client := github.GetGitHubOAuthConfig().Client(r.Context(), token)
	response, err := client.Get("https://api.github.com/user")
	if err != nil {
		http.Error(w, "❌ Failed to fetch user: "+err.Error(), http.StatusInternalServerError)
	}

	defer response.Body.Close()

	var githubUser struct {
		Login string `json:"login"`
		Email string `json:"email"`
	}

	json.NewDecoder(response.Body).Decode(&githubUser)

	if githubUser.Login == "" {
		http.Error(w, "❌ Invalid GitHub response", http.StatusInternalServerError)
		return
	}

	dbInstance := config.DB

	// To check if user exists
	var existingUser models.User
	err = dbInstance.Where("git_hub_username = ?", githubUser.Login).First(&existingUser).Error

	if err == gorm.ErrRecordNotFound {
		// To create new user
		newUser := models.User{
			Email:          githubUser.Email,
			GitHubUsername: githubUser.Login,
			GitHubToken:    token.AccessToken,
		}

		if err := dbInstance.Create(&newUser).Error; err != nil {
			http.Error(w, "❌ Failed to create user", http.StatusInternalServerError)
			return
		}

		// fmt.Fprintf(w, "New User Created: , %s!", githubUser.Login)
	}

	// To redirect to the frontend dashboard
	frontendURL := os.Getenv("FRONTEND_URL")

	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}

	redirectURL := fmt.Sprintf("%s/dashboard", frontendURL)

	fmt.Printf("Redirecting user to: %s\n", redirectURL)
	http.Redirect(w, r, redirectURL, http.StatusFound)

}
