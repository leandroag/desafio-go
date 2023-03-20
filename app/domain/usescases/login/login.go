package login

import (
	"context"

	"github.com/leandroag/desafio/app/domain/entities"
	"golang.org/x/crypto/bcrypt"
)

type accountRepository interface {
	GetAccountByCPF(ctx context.Context, cpf string) (entities.Account, error)
}

type cryptService interface {
	GenerateToken(accountID int32) (string, error)
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

func (service loginService) Authenticate(ctx context.Context, cpf string, secret string) (string, error) {
	account, err := service.accountRepository.GetAccountByCPF(ctx, cpf)
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
