package login

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/leandroag/desafio/app/domain/entities"
)

type loginService interface {
	Authenticate(ctx context.Context, cpf string, secret string) (string, error)
}

type LoginHandler struct {
	loginUseCase loginService
}

func NewLoginHandler(loginUseCase loginService) *LoginHandler {
	return &LoginHandler{
		loginUseCase,
	}
}

func (handler *LoginHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", handler.login).Methods(http.MethodPost)
}

func (handler *LoginHandler) login(w http.ResponseWriter, r *http.Request) {
	var login entities.Login

	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := handler.loginUseCase.Authenticate(r.Context(), login.CPF, login.Secret)
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
