package transfer

import (
	"errors"

	"github.com/leandroag/desafio/app/domain/entities"
)

type transferRepository interface {
	CreateTransfer(transfer entities.Transfer) error
	GetTransfersByAccountID(accountID string) ([]entities.Transfer, error)
}

type accountRepository interface {
	GetAccountByID(accountID string) (entities.Account, error)
	GetAccountByCPF(CPF string) (entities.Account, error)
	UpdateAccount(account entities.Account) error
}

type cryptService interface {
	GetAccountByToken(token string) (string, error)
}

type transferService struct {
	accountRepository  accountRepository
	transferRepository transferRepository
	cryptService       cryptService
}

func NewTransferService(accountRepository accountRepository, transferRepository transferRepository, cryptService cryptService) *transferService {
	return &transferService{
		accountRepository,
		transferRepository,
		cryptService,
	}
}

func (service *transferService) CreateTransfer(token string, transfer entities.Transfer) error {
	// Implementar a função de buscar informações da conta a partir do token do usuário autenticado atualmente
	accountOriginID, err := service.cryptService.GetAccountByToken(token)
	if err != nil {
		return err
	}

	if accountOriginID != transfer.AccountOriginID {
		return errors.New("Conta de Origem inválida")
	}

	accountOrigin, err := service.accountRepository.GetAccountByID(accountOriginID)
	if err != nil {
		return err
	}

	accountDestination, err := service.accountRepository.GetAccountByID(transfer.AccountDestinationID)
	if err != nil {
		return err
	}

	if accountOrigin.Balance < transfer.Amount {
		return errors.New("O balanço da Conta de Origem é insuficiente")
	}

	accountOrigin.Balance -= transfer.Amount
	accountDestination.Balance += transfer.Amount

	// Atualiza o balanço da conta de origem
	if err := service.accountRepository.UpdateAccount(accountOrigin); err != nil {
		return err
	}

	// Atualiza o balanço da conta de destino
	if err := service.accountRepository.UpdateAccount(accountDestination); err != nil {
		return err
	}

	return service.transferRepository.CreateTransfer(transfer)
}

func (service *transferService) GetTransfersByAccountID(accountID string) ([]entities.Transfer, error) {
	return service.transferRepository.GetTransfersByAccountID(accountID)
}
