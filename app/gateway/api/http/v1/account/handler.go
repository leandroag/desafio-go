package account

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/leandroag/desafio/app/domain/entities"
	"github.com/leandroag/desafio/app/dtos"
)

type accountService interface {
	CreateAccount(ctx context.Context, account dtos.AccountDTO) error
	GetAccountBalance(ctx context.Context, accountID int32) (float64, error)
	GetAccounts(ctx context.Context) ([]entities.Account, error)
}

type AccountHandler struct {
	accountUseCase accountService
}

func NewAccountHandler(accountUseCase accountService) *AccountHandler {
	return &AccountHandler{
		accountUseCase,
	}
}

func (handler AccountHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/accounts", handler.createAccount).Methods(http.MethodPost)
	router.HandleFunc("/accounts", handler.listAccounts).Methods(http.MethodGet)
	router.HandleFunc("/accounts/{account_id}/balance", handler.getAccountBalance).Methods(http.MethodGet)
}

func (handler AccountHandler) createAccount(w http.ResponseWriter, r *http.Request) {
	var account dtos.AccountDTO

	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = handler.accountUseCase.CreateAccount(r.Context(), account)
	if err != nil {
		http.Error(w, "Error creating account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (handler *AccountHandler) listAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := handler.accountUseCase.GetAccounts(r.Context())
	if err != nil {
		http.Error(w, "Error getting accounts", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(accounts)
}

func (handler *AccountHandler) getAccountBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	accountID, err := strconv.ParseInt(vars["account_id"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	balance, err := handler.accountUseCase.GetAccountBalance(r.Context(), int32(accountID))
	if err != nil {
		http.Error(w, "Error getting account balance", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]float64{"balance": balance})
}
