package api

import (
	"fmt"
	"github.com/cloudmusic-dev/backend/database"
	"github.com/gorilla/mux"
	"net/http"
)

func handleListPlaylists(w http.ResponseWriter, r *http.Request) {
	var playlists []database.Playlist
	database.DB.Find(&playlists)

	w.WriteHeader(http.StatusNotImplemented)
	fmt.Fprintf(w, "%d", len(playlists))
}

func CreatePlaylistRouter(router *mux.Router) {
	router.HandleFunc("/playlists", handleListPlaylists)
}
