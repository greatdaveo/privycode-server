package routes

import (
	"net/http"

	"github.com/greatdaveo/privycode-server/internal/handlers"
)

func APIRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to PrivyCode 👋"))
	})

	mux.HandleFunc("/github/login", handlers.GitHubLoginHandler)
	mux.HandleFunc("/github/callback", handlers.GitHubCallbackHandler)

	mux.HandleFunc("/generate-viewer-link", handlers.GenerateViewerLinkHandler)

	mux.HandleFunc("/view/", handlers.ViewerAccessHandler)
	
}
