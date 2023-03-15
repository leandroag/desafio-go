package transfer

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/leandroag/desafio/app/domain/entities"
	"github.com/leandroag/desafio/app/domain/usescases/transfer"
)

type TransferDTO struct {
	AccountOriginID      uint64 `json:"account_origin_id"`
	AccountDestinationID uint64 `json:"account_destination_id"`
	Amount               uint64 `json:"amount"`
}

type cryptService interface {
	GetAccountByToken(token string) (string, error)
	GetToken(accountID string) string 
}

type TransferHandler struct {
	transferUseCase transfer.TransferService
	cryptService    cryptService
}

func NewTransferHandler(transferUseCase transfer.TransferService, cryptService cryptService) *TransferHandler {
	return &TransferHandler{
		transferUseCase,
		cryptService,
	}
}

func (handler TransferHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/transfers ", handler.getTransfers).Methods(http.MethodGet)
	router.HandleFunc("/transfers ", handler.createTransfer).Methods(http.MethodPost)
}

// Handler para a rota GET /transfers
func (handler TransferHandler) getTransfers(w http.ResponseWriter, r *http.Request) {
	// Recupera o valor do cabeçalho "Authorization" da requisição
	authHeader := r.Header.Get("Authorization")

	// Separa o valor do token da string "Bearer "
	token := strings.Split(authHeader, "Bearer ")[1]

	// Busca informações da conta a partir do token do usuário autenticado atualmente
	accountOriginID, err := handler.cryptService.GetAccountByToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Recupera a lista de transferências da conta de origem
	transferList, err := handler.transferUseCase.GetTransfersByAccountID(accountOriginID)
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
	var transfer entities.Transfer

	err := json.NewDecoder(r.Body).Decode(&transfer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Recupera o valor do cabeçalho "Authorization" da requisição
	authHeader := r.Header.Get("Authorization")

	// Separa o valor do token da string "Bearer "
	token := strings.Split(authHeader, "Bearer ")[1]

	// Realiza a transferência
	err = handler.transferUseCase.CreateTransfer(token, transfer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
