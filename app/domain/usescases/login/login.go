package login

import (
	"github.com/leandroag/desafio/app/domain/entities"
	"golang.org/x/crypto/bcrypt"
)

type accountRepository interface {
	GetAccountByCPF(cpf string) (entities.Account, error)
}

type cryptService interface {
	GenerateToken(accountID string) (string, error)
}

type LoginService interface {
	Authenticate(cpf string, secret string) (string, error)
}

type loginService struct {
	accountRepository accountRepository
	cryptService      cryptService
}

func NewLoginService(accountRepository accountRepository, cryptService cryptService) *loginService {
	return &loginService{
		accountRepository,
		cryptService,
	}
}

func (service loginService) Authenticate(cpf string, secret string) (string, error) {
	account, err := service.accountRepository.GetAccountByCPF(cpf)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Secret), []byte(secret))
	if err != nil {
		// a senha fornecida é inválida, fazer tratamento para retornar erro
	}

	// Fazer a geração do token
	token, err := service.cryptService.GenerateToken(account.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
