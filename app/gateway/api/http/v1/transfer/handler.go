package transfer

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/leandroag/desafio/app/domain/entities"
	"github.com/leandroag/desafio/app/dtos"
)

type transferService interface {
	CreateTransfer(ctx context.Context, token string, transfer dtos.TransferDTO) error
	GetTransfersByAccountID(ctx context.Context, accountID int32) ([]entities.Transfer, error)
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

func (handler TransferHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/transfers", handler.getTransfers).Methods(http.MethodGet)
	router.HandleFunc("/transfers", handler.createTransfer).Methods(http.MethodPost)
}

// Handler para a rota GET /transfers
func (handler TransferHandler) getTransfers(w http.ResponseWriter, r *http.Request) {
	// Recupera o valor do cabeçalho "Authorization" da requisição
	token := r.Header.Get("Authorization")

	// Busca informações da conta a partir do token do usuário autenticado atualmente
	accountOriginID, err := handler.cryptService.GetAccountByToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Recupera a lista de transferências da conta de origem
	transferList, err := handler.transferUseCase.GetTransfersByAccountID(r.Context(), accountOriginID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Serializa a lista de transferências em formato JSON e envia na resposta
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transferList)
}

// Handler para a rota POST /transfers
func (handler TransferHandler) createTransfer(w http.ResponseWriter, r *http.Request) {
	var transfer dtos.TransferDTO

	err := json.NewDecoder(r.Body).Decode(&transfer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Recupera o valor do cabeçalho "Authorization" da requisição
	token := r.Header.Get("Authorization")

	// Realiza a transferência
	err = handler.transferUseCase.CreateTransfer(r.Context(), token, transfer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
