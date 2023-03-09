package account

import (
	"github.com/leandroag/desafio/app/domain/entities"
)

type accountRepository interface {
	CreateAccount(account entities.Account) error
	GetAccountBalance(accountID string) (float64, error)
	GetAllAccounts() ([]entities.Account, error)
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

func (service *accountService) GetAccounts() ([]entities.Account, error) {
	return service.accountRepository.GetAllAccounts()
}

func (service *accountService) GetAccountBalance(accountID string) (float64, error) {
	return service.accountRepository.GetAccountBalance(accountID)
}

func (service *accountService) CreateAccount(account entities.Account) error {
	passwordHash, err := service.cryptService.HashSecret(account.Secret)
	if err != nil {
		return err
	}

	account.Secret = passwordHash

	return service.accountRepository.CreateAccount(account)
}
