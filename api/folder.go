package api

import (
	"encoding/json"
	"github.com/cloudmusic-dev/backend/authorization"
	"github.com/cloudmusic-dev/backend/database"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type Folder struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	ParentFolder string `json:"parentFolder"`
}

func databaseFolderToApi(folder database.Folder) Folder {
	return Folder{
		Id:           folder.ID.String(),
		Name:         folder.Name,
		ParentFolder: convertUUIDToString(folder.ParentFolder),
	}
}

func handleGetFolders(w http.ResponseWriter, r *http.Request) {
	authorized, userId := authorization.ValidateRequest(r)
	if !authorized {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var folders []database.Folder
	database.DB.Find(&folders, "owner = ?", userId)

	ret := make([]Folder, len(folders))
	for i := 0; i < len(ret); i++ {
		ret[i] = databaseFolderToApi(folders[i])
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&ret); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleCreateFolder(w http.ResponseWriter, r *http.Request) {
	authorized, userId := authorization.ValidateRequest(r)
	if !authorized {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req Folder
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	folder := database.Folder{
		Owner:        *userId,
		Name:         req.Name,
		ParentFolder: convertStringToUUID(req.ParentFolder),
		CreatedAt:    time.Now(),
	}
	if err := database.DB.Save(&folder).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(databaseFolderToApi(folder)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func CreateFolderRouter(router *mux.Router) {
	router.HandleFunc("/folders", handleGetFolders).Methods("GET")
	router.HandleFunc("/folder", handleCreateFolder).Methods("POST")
}
