package transfer

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/leandroag/desafio/app/dtos"
)

type transferService interface {
	CreateTransfer(ctx context.Context, token string, transfer dtos.TransferDTO) error
	GetTransfersByAccountID(ctx context.Context, accountID int32) ([]dtos.TransferDTO, error)
}

type cryptService interface {
	GetAccountByToken(token string) (int32, error)
}

type TransferHandler struct {
	transferUseCase transferService
	cryptService    cryptService
}

func NewTransferHandler(transferUseCase transferService, cryptService cryptService) *TransferHandler {
	return &TransferHandler{
		transferUseCase,
		cryptService,
	}
}

func (h TransferHandler) RegisterRoutes(router *chi.Mux) {
	router.Get("/transfers", h.getTransfers)
	router.Post("/transfers", h.createTransfer)
}

// getTransfers retrieves a list of transfers made by the authenticated account
// @Summary Retrieves a list of transfers
// @Description Retrieves a list of transfers made by the authenticated account.
// @Tags Transfers
// @Security BearerAuth
// @Success 200 {array} dtos.TransferDTO
// @Failure 400 {string} string "Invalid request payload"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /transfers [get]
func (h TransferHandler) getTransfers(w http.ResponseWriter, r *http.Request) {
	// Recupera o valor do cabeçalho "Authorization" da requisição
	token := r.Header.Get("Authorization")

	// Busca informações da conta a partir do token do usuário autenticado atualmente
	accountOriginID, err := h.cryptService.GetAccountByToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Recupera a lista de transferências da conta de origem
	transferList, err := h.transferUseCase.GetTransfersByAccountID(r.Context(), accountOriginID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Serializa a lista de transferências em formato JSON e envia na resposta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transferList)
}

// createTransfer creates a transfer between two accounts.
// @Summary Create transfer
// @Description Creates a transfer between two accounts.
// @Tags Transfers
// @Accept json
// @Produce json
// @Param transfer body dtos.TransferDTO true "Transfer object"
// @Success 201 {object} string "Transfer created successfully"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Error creating transfer"
// @Router /transfers [post]
func (h TransferHandler) createTransfer(w http.ResponseWriter, r *http.Request) {
	var transfer dtos.TransferDTO

	err := json.NewDecoder(r.Body).Decode(&transfer)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Recupera o valor do cabeçalho "Authorization" da requisição
	token := r.Header.Get("Authorization")

	// Realiza a transferência
	err = h.transferUseCase.CreateTransfer(r.Context(), token, transfer)
	if err != nil {
		http.Error(w, "Error creating transfer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
