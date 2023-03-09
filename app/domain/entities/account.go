package entities

import "time"

// Account representa uma conta bancária
type Account struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CPF       string    `json:"cpf"`
	Secret    string    `json:"-"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}
