package dtos

type ListAccountDTO struct {
	Name    string  `json:"name"`
	CPF     string  `json:"cpf"`
	Balance float64 `json:"balance"`
}
