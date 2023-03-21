package account

import (
	"context"

	"github.com/leandroag/desafio/app/domain/entities"
	"github.com/leandroag/desafio/app/dtos"
)

type accountRepository interface {
	CreateAccount(ctx context.Context, account entities.Account) error
	GetAccountBalance(ctx context.Context, accountID int32) (float64, error)
	GetAllAccounts(ctx context.Context) ([]entities.Account, error)
}

type cryptService interface {
	HashSecret(secret string) (string, error)
}

type accountService struct {
	accountRepository accountRepository
	cryptService      cryptService
}

func NewAccountService(accountRepository accountRepository, cryptService cryptService) *accountService {
	return &accountService{
		accountRepository,
		cryptService,
	}
}

func (s accountService) GetAccounts(ctx context.Context) ([]dtos.ListAccountDTO, error) {
	allAccounts, err := s.accountRepository.GetAllAccounts(ctx)

	if err != nil {
		return []dtos.ListAccountDTO{}, err
	}

	var accountsList []dtos.ListAccountDTO

	for _, account := range allAccounts {
		accountDTO := dtos.ListAccountDTO{
			Name:    account.Name,
			CPF:     account.CPF,
			Balance: account.Balance,
		}

		accountsList = append(accountsList, accountDTO)
	}

	return accountsList, nil
}

func (s accountService) GetAccountBalance(ctx context.Context, accountID int32) (float64, error) {
	return s.accountRepository.GetAccountBalance(ctx, accountID)
}

func (s accountService) CreateAccount(ctx context.Context, account dtos.CreateAccountDTO) error {
	passwordHash, err := s.cryptService.HashSecret(account.Secret)
	if err != nil {
		return err
	}

	account.Secret = passwordHash

	accountToSave := account.ToCreateAccountDomain()

	return s.accountRepository.CreateAccount(ctx, *accountToSave)
}
