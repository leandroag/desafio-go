package entities

import "time"

// Transfer representa uma transferência de dinheiro entre contas bancárias
type Transfer struct {
	ID                   string
	AccountOriginID      int32
	AccountDestinationID int32
	Amount               float64
	CreatedAt            time.Time
}
