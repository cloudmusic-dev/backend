package api

import (
	"encoding/json"
	"github.com/cloudmusic-dev/backend/authorization"
	"github.com/cloudmusic-dev/backend/database"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type Playlist struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	FolderId  string    `json:"folderId"`
	CreatedAt time.Time `json:"createdAt"`
}

func databasePlaylistToApi(playlist database.Playlist) Playlist {
	ret := Playlist{
		ID:        playlist.ID.String(),
		Name:      playlist.Name,
		CreatedAt: playlist.CreatedAt,
	}

	if playlist.ParentFolder != nil {
		ret.FolderId = playlist.ParentFolder.String()
	}

	return ret
}

func handleListPlaylists(w http.ResponseWriter, r *http.Request) {
	authorized, userId := authorization.ValidateRequest(r)
	if !authorized {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var playlists []database.Playlist
	database.DB.Where("owner = ?", userId).Find(&playlists)

	ret := make([]Playlist, len(playlists))
	for i := 0; i < len(ret); i++ {
		ret[i] = databasePlaylistToApi(playlists[i])
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(ret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Failed to encode response: %v", err)
	}
}

func handleCreatePlaylist(w http.ResponseWriter, r *http.Request) {
	authorized, userId := authorization.ValidateRequest(r)
	if !authorized {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var request Playlist
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	playlist := database.Playlist{
		Name: request.Name,
		Owner: *userId,
		CreatedAt: time.Now(),
	}
	if folderId, err := uuid.Parse(request.FolderId); err != nil && request.FolderId != "" {
		playlist.ParentFolder = &folderId
	}
	err = database.DB.Create(&playlist).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(databasePlaylistToApi(playlist))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleDeletePlaylist(w http.ResponseWriter, r *http.Request) {
	authorized, userId := authorization.ValidateRequest(r)
	if !authorized {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)

	var playlist database.Playlist
	if database.DB.First(&playlist, "ID = ?", vars["id"]).RecordNotFound() || playlist.Owner != *userId {
		http.Error(w, "playlist not found", http.StatusNotFound)
		return
	}

	if err := database.DB.Delete(&playlist).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleUpdatePlaylist(w http.ResponseWriter, r *http.Request) {
	authorized, userId := authorization.ValidateRequest(r)
	if !authorized {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)

	var playlist database.Playlist
	if database.DB.First(&playlist, "ID = ?", vars["id"]).RecordNotFound() || playlist.Owner != *userId {
		http.Error(w, "playlist not found", http.StatusNotFound)
		return
	}

	var request Playlist
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	playlist.Name = request.Name
	if folderId, err := uuid.Parse(request.FolderId); err != nil {
		if request.FolderId != "" {
			playlist.ParentFolder = &folderId
		} else {
			playlist.ParentFolder = nil
		}
	}

	if err := database.DB.Save(&playlist).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(databasePlaylistToApi(playlist))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func CreatePlaylistRouter(router *mux.Router) {
	router.HandleFunc("/playlists", handleListPlaylists).Methods("GET")
	router.HandleFunc("/playlist", handleCreatePlaylist).Methods("POST")
	router.HandleFunc("/playlist/{id}", handleDeletePlaylist).Methods("DELETE")
	router.HandleFunc("/playlist/{id}", handleUpdatePlaylist).Methods("PATCH")
}
