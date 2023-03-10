package account

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/leandroag/desafio/app/domain/entities"
	"github.com/leandroag/desafio/app/domain/usescases/account"
)

type AccountHandler struct {
	accountUseCase account.AccountService
}

func NewAccountHandler(accountUseCase account.AccountService) *AccountHandler {
	return &AccountHandler{
		accountUseCase,
	}
}

func (handler *AccountHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/accounts", handler.createAccount).Methods(http.MethodPost)
	router.HandleFunc("/accounts", handler.listAccounts).Methods(http.MethodGet)
	router.HandleFunc("/accounts/{account_id}/balance", handler.getAccountBalance).Methods(http.MethodGet)
}

func (handler *AccountHandler) createAccount(w http.ResponseWriter, r *http.Request) {
	var account entities.Account

	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = handler.accountUseCase.CreateAccount(account)
	if err != nil {
		http.Error(w, "Error creating account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (handler *AccountHandler) listAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := handler.accountUseCase.GetAccounts()
	if err != nil {
		http.Error(w, "Error getting accounts", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(accounts)
}

func (handler *AccountHandler) getAccountBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	accountID := vars["account_id"]

	balance, err := handler.accountUseCase.GetAccountBalance(accountID)
	if err != nil {
		http.Error(w, "Error getting account balance", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]float64{"balance": balance})
}
