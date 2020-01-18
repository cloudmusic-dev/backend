package api

import (
	"encoding/json"
	"github.com/cloudmusic-dev/backend/authorization"
	"github.com/cloudmusic-dev/backend/database"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PublicUser struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}

type LoginResponse struct {
	ApiKey string     `json:"apiKey"`
	User   PublicUser `json:"user"`
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var user database.User
	if database.DB.Where("username = ?", req.Username).First(&user).RecordNotFound() {
		// Do bcrypt check anyways, this is to prevent to do timing attacks if the user exists or not
		bcrypt.CompareHashAndPassword([]byte(""), []byte(req.Password))

		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err == nil {
		// Password is correct, generate access token
		token, err := authorization.CreateToken(user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Failed to generate jwt token: %v", err)
			return
		}

		response := LoginResponse{
			ApiKey: token,
			User: PublicUser{
				ID:          user.ID.String(),
				DisplayName: user.Username,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("Failed to encode response: %v", err)
		}
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

func CreateAccountRouter(router *mux.Router) {
	router.HandleFunc("/login", handleLogin)
}
