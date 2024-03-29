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

// createAccount creates a new account.
// @Summary Create account
// @Description Creates a new account.
// @Tags Account
// @Accept json
// @Produce json
// @Param account body dtos.CreateAccountDTO true "Account object"
// @Success 201 {object} string "Account created successfully"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Error creating account"
// @Router /accounts [post]
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

// listAccounts returns a list of accounts.
// @Summary List accounts
// @Description Returns a list of accounts.
// @Tags Account
// @Accept json
// @Produce json
// @Success 200 {array} dtos.ListAccountDTO 
// @Failure 500 {object} string "Error getting accounts"
// @Router /accounts [get]
func (h *AccountHandler) listAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.accountUseCase.GetAccounts(r.Context())
	if err != nil {
		http.Error(w, "Error getting accounts", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(accounts)
}

// getAccountBalance returns the balance of an account.
// @Summary Get account balance
// @Description Returns the balance of an account.
// @Tags Account
// @Accept json
// @Produce json
// @Param account_id path int true "Account ID"
// @Success 200 {object} float64
// @Failure 400 {object} string "Invalid account ID"
// @Failure 500 {object} string "Error getting account balance"
// @Router /accounts/{account_id}/balance [get]
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
