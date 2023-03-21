package entities

import "time"

// Account representa uma conta bancária
type Account struct {
	ID        int32
	Name      string
	CPF       string
	Secret    string
	Balance   float64
	CreatedAt time.Time
}
