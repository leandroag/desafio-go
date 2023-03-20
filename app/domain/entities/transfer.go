package entities

import "time"

// Transfer representa uma transferência de dinheiro entre contas bancárias
type Transfer struct {
	ID                   string    `json:"id"`
	AccountOriginID      int32     `json:"account_origin_id"`
	AccountDestinationID int32     `json:"account_destination_id"`
	Amount               float64   `json:"amount"`
	CreatedAt            time.Time `json:"created_at"`
}
