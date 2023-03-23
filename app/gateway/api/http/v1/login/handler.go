package login

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/leandroag/desafio/app/dtos"
)

type loginService interface {
	Authenticate(ctx context.Context, Login dtos.LoginDTO) (string, error)
}

type LoginHandler struct {
	loginUseCase loginService
}

func NewLoginHandler(loginUseCase loginService) *LoginHandler {
	return &LoginHandler{
		loginUseCase,
	}
}

func (h *LoginHandler) RegisterRoutes(router *chi.Mux) {
	router.Post("/login", h.login)
}

func (h *LoginHandler) login(w http.ResponseWriter, r *http.Request) {
	var login dtos.LoginDTO

	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	token, err := h.loginUseCase.Authenticate(r.Context(), login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response := token
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
