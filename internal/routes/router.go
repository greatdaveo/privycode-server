package routes

import "net/http"

func APIRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to PrivyCode 👋"))
	})

}
