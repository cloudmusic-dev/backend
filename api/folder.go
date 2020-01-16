package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

func handleGetFolders(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func CreateFolderRouter(router *mux.Router) {
	router.HandleFunc("/folders", handleGetFolders)
}
