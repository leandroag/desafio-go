package login

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/leandroag/desafio/app/domain/entities"
	"github.com/leandroag/desafio/app/domain/usescases/login"
)

type LoginHandler struct {
	loginUseCase login.LoginService
}

func NewLoginHandler(loginUseCase login.LoginService) *LoginHandler {
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

	token, err := handler.loginUseCase.Authenticate(login.CPF, login.Secret)
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
