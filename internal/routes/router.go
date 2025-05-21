package routes

import (
	"net/http"

	"github.com/greatdaveo/privycode-server/internal/handlers"
	"github.com/greatdaveo/privycode-server/internal/middleware"
)

func APIRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to PrivyCode 👋"))
	})

	mux.HandleFunc("/github/login", handlers.GitHubLoginHandler)
	mux.HandleFunc("/dashboard", middleware.AuthMiddleware(handlers.DashboardHandler))
	mux.HandleFunc("/github/callback", handlers.GitHubCallbackHandler)

	mux.HandleFunc("/generate-viewer-link", middleware.AuthMiddleware(handlers.GenerateViewerLinkHandler))

	mux.HandleFunc("/view/", handlers.ViewerAccessHandler)
	mux.HandleFunc("/view-files/", handlers.ViewFileHandler)
	mux.HandleFunc("/view-folder/", handlers.ViewerFolderHandler)

	mux.HandleFunc("/view-info/", handlers.ViewUserInfoHandler)
}
