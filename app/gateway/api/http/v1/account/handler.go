package account

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/leandroag/desafio/app/dtos"
)

type accountService interface {
	CreateAccount(ctx context.Context, account dtos.CreateAccountDTO) error
	GetAccountBalance(ctx context.Context, accountID int32) (float64, error)
	GetAccounts(ctx context.Context) ([]dtos.ListAccountDTO, error)
}

type AccountHandler struct {
	accountUseCase accountService
}

func NewAccountHandler(accountUseCase accountService) *AccountHandler {
	return &AccountHandler{
		accountUseCase,
	}
}

func (h AccountHandler) RegisterRoutes(router *chi.Mux) {
	router.Post("/accounts", h.createAccount)
	router.Get("/accounts", h.listAccounts)
	router.Get("/accounts/{account_id}/balance", h.getAccountBalance)
}

func (h AccountHandler) createAccount(w http.ResponseWriter, r *http.Request) {
	var account dtos.CreateAccountDTO

	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = h.accountUseCase.CreateAccount(r.Context(), account)
	if err != nil {
		http.Error(w, "Error creating account", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AccountHandler) listAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.accountUseCase.GetAccounts(r.Context())
	if err != nil {
		http.Error(w, "Error getting accounts", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(accounts)
}

func (h *AccountHandler) getAccountBalance(w http.ResponseWriter, r *http.Request) {
	accountIDString := chi.URLParam(r, "account_id")

	accountID, err := strconv.ParseInt(accountIDString, 10, 32)

	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	balance, err := h.accountUseCase.GetAccountBalance(r.Context(), int32(accountID))
	if err != nil {
		http.Error(w, "Error getting account balance", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]float64{"balance": balance})
}
