package github

import (
	"context"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func GetGitHubOAuthConfig() *oauth2.Config {

	return &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Endpoint:     github.Endpoint,
		RedirectURL:  os.Getenv("GITHUB_CALLBACK_URL"),
		Scopes:       []string{"read:user", "repo"},
	}
}

// To generate GitHub OAuth URL
func GetAuthURL(state string) string {
	return GetGitHubOAuthConfig().AuthCodeURL(state)
}

// To exchange code for access token
func ExchangeCodeForToken(code string) (*oauth2.Token, error) {
	ctx := context.Background()
	return GetGitHubOAuthConfig().Exchange(ctx, code)
}
