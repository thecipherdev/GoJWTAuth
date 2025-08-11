package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/thecipherdev/goauth/dto"
	"github.com/thecipherdev/goauth/mock"
	"github.com/thecipherdev/goauth/model"
	"github.com/thecipherdev/goauth/utils"
)

type UserHandler struct{}

const DummyHash = "$argon2id$v=19$m=65536,t=3,p=4$MTIzNDU2Nzg5MGFiY2RlZg$s+EbQVuzvGQiT+np23+6ouYZdELuPcUFKLZ/nxEnIgQ"

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) UserRouter(router *http.ServeMux) {
	router.HandleFunc("POST /register", handleRegister)
	router.HandleFunc("POST /login", handleLogin)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	var payload dto.UserLogin
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := utils.GetByUsername(payload.Username)

	if err != nil {
		// Verify dummy hash for Timing Equalization
		utils.VerifyPassword(payload.Password, DummyHash)
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	ok, err := utils.VerifyPassword(payload.Password, user.Password)

	if err != nil || !ok {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateAccessToken()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to genereate access token: %v", err), http.StatusInternalServerError)
		return
	}

	refresh, err := utils.GenerateRefreshToken()

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate refresh token: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Println("Is User Logged in:", ok)

	w.Header().Set("Content-Type", "application/json")
	data := map[string]string{
		"access_token":  token,
		"refresh_token": refresh,
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}

}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	var payload dto.UserDTO

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body %v", err), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	enc, err := utils.HashPassword(payload.Password)

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed while hashing %v", err), http.StatusInternalServerError)
		return
	}

	mock.Users = append(mock.Users, model.User{
		Username: payload.Username,
		Password: enc,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	data := map[string]any{
		"message": "user created.",
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}

}
