package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func handleListPlaylists(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "500")
}

func CreatePlaylistRouter(router *mux.Router) {
	router.HandleFunc("/playlists", handleListPlaylists)
}
