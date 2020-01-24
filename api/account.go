package api

import (
	"database/sql"
	"encoding/json"
	"github.com/cloudmusic-dev/backend/authorization"
	"github.com/cloudmusic-dev/backend/database"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
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

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type RegisterErrorResponse struct {
	Error string `json:"error"`
}

type RegisterSuccessResponse struct {
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

func sendRegisterError(error string, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(RegisterErrorResponse{Error: error}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// todo validate input

	// Check if username is already in use
	var user database.User
	if !database.DB.First(&user, "username = ?", req.Username).RecordNotFound() {
		sendRegisterError("Username is already in use", w)
		return
	}

	// Check if email is already in use
	if !database.DB.First(&user, "password = ?", req.Password).RecordNotFound() {
		sendRegisterError("Email is already in use", w)
		return
	}

	// Create password hash
	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Store user to database
	user = database.User{
		Username:       req.Username,
		Email:          req.Email,
		Password:       string(password),
		Activated:      sql.NullBool{Bool: true, Valid: true},
		ActivationCode: "",
		CreatedAt:      time.Now(),
	}
	if err := database.DB.Save(&user).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create jwt token for client
	token, err := authorization.CreateToken(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send token and user profile to client
	if err := json.NewEncoder(w).Encode(RegisterSuccessResponse{
		ApiKey: token,
		User: PublicUser{
			ID:          user.ID.String(),
			DisplayName: user.Username,
		},
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func CreateAccountRouter(router *mux.Router) {
	router.HandleFunc("/login", handleLogin)
	router.HandleFunc("/register", handleRegister)
}
