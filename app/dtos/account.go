package dtos

import "github.com/leandroag/desafio/app/domain/entities"

type AccountDTO struct {
	Name    string  `json:"name"`
	CPF     string  `json:"cpf"`
	Secret  string  `json:"secret"`
	Balance float64 `json:"balance"`
}

func (dto *AccountDTO) ToAccountDomain() *entities.Account {
	return &entities.Account{
		Name:    dto.Name,
		CPF:     dto.CPF,
		Secret:  dto.Secret,
		Balance: dto.Balance,
	}
}
