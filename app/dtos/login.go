package dtos

import "github.com/leandroag/desafio/app/domain/entities"

type LoginDTO struct {
	CPF    string `json:"cpf"`
	Secret string `json:"secret"`
}

func (dto *LoginDTO) ToLoginDomain() *entities.Login {
	return &entities.Login{
		CPF:    dto.CPF,
		Secret: dto.Secret,
	}
}
