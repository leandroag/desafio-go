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

func (s transferService) CreateTransfer(ctx context.Context, token string, transfer dtos.TransferDTO) error {
	// Busca informações da conta a partir do token do usuário autenticado atualmente
	accountOriginID, err := s.cryptService.GetAccountByToken(token)
	if err != nil {
		return err
	}

	if accountOriginID != transfer.AccountOriginID {
		return errors.New("conta de origem inválida")
	}

	accountOrigin, err := s.accountRepository.GetAccountByID(ctx, accountOriginID)
	if err != nil {
		return err
	}

	accountDestination, err := s.accountRepository.GetAccountByID(ctx, transfer.AccountDestinationID)
	if err != nil {
		return err
	}

	if accountOrigin.Balance < transfer.Amount {
		return errors.New("o balanço da conta de origem é insuficiente")
	}

	accountOrigin.Balance -= transfer.Amount
	accountDestination.Balance += transfer.Amount

	// Atualiza o balanço da conta de origem
	if err := s.accountRepository.UpdateAccount(ctx, accountOrigin); err != nil {
		return err
	}

	// Atualiza o balanço da conta de destino
	if err := s.accountRepository.UpdateAccount(ctx, accountDestination); err != nil {
		return err
	}

	transferToSave := transfer.ToTransferDomain()

	return s.transferRepository.CreateTransfer(ctx, *transferToSave)
}

func (s transferService) GetTransfersByAccountID(ctx context.Context, accountID int32) ([]dtos.TransferDTO, error) {
	transfersByAccount, err := s.transferRepository.GetTransfersByAccountID(ctx, accountID)

	if err != nil {
		return []dtos.TransferDTO{}, err
	}

	var transfersList []dtos.TransferDTO

	for _, transfer := range transfersByAccount {
		transferDTO := dtos.TransferDTO{
			AccountOriginID:      transfer.AccountOriginID,
			AccountDestinationID: transfer.AccountDestinationID,
			Amount:               transfer.Amount,
		}

		transfersList = append(transfersList, transferDTO)
	}

	return transfersList, nil
}
