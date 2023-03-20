package transfer

import (
	"context"
	"errors"

	"github.com/leandroag/desafio/app/domain/entities"
	"github.com/leandroag/desafio/app/dtos"
)

type transferRepository interface {
	CreateTransfer(ctx context.Context, transfer entities.Transfer) error
	GetTransfersByAccountID(ctx context.Context, accountID int32) ([]entities.Transfer, error)
}

type accountRepository interface {
	GetAccountByID(ctx context.Context, accountID int32) (entities.Account, error)
	GetAccountByCPF(ctx context.Context, CPF string) (entities.Account, error)
	UpdateAccount(ctx context.Context, account entities.Account) error
}

type cryptService interface {
	GetAccountByToken(token string) (int32, error)
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

func (service transferService) CreateTransfer(ctx context.Context, token string, transfer dtos.TransferDTO) error {
	// Busca informações da conta a partir do token do usuário autenticado atualmente
	accountOriginID, err := service.cryptService.GetAccountByToken(token)
	if err != nil {
		return err
	}

	if accountOriginID != transfer.AccountOriginID {
		return errors.New("conta de origem inválida")
	}

	accountOrigin, err := service.accountRepository.GetAccountByID(ctx, accountOriginID)
	if err != nil {
		return err
	}

	accountDestination, err := service.accountRepository.GetAccountByID(ctx, transfer.AccountDestinationID)
	if err != nil {
		return err
	}

	if accountOrigin.Balance < transfer.Amount {
		return errors.New("o balanço da conta de origem é insuficiente")
	}

	accountOrigin.Balance -= transfer.Amount
	accountDestination.Balance += transfer.Amount

	// Atualiza o balanço da conta de origem
	if err := service.accountRepository.UpdateAccount(ctx, accountOrigin); err != nil {
		return err
	}

	// Atualiza o balanço da conta de destino
	if err := service.accountRepository.UpdateAccount(ctx, accountDestination); err != nil {
		return err
	}

	transferToSave := transfer.ToTransferDomain()

	return service.transferRepository.CreateTransfer(ctx, *transferToSave)
}

func (service transferService) GetTransfersByAccountID(ctx context.Context, accountID int32) ([]entities.Transfer, error) {
	return service.transferRepository.GetTransfersByAccountID(ctx, accountID)
}
