package postgres

import (
	"context"
	"time"

	"github.com/leandroag/desafio/app/domain/entities"
)

type transferRepository struct {
	Querier
}

func NewTransferRepository(querier Querier) *transferRepository {
	return &transferRepository{
		querier,
	}
}

func (r transferRepository) CreateTransfer(ctx context.Context, accountOriginID string, accountDestinationID string, amount float64) (entities.Transfer, error) {
	const query = "INSERT INTO transfers(account_origin_id, account_destination_id, amount) VALUES($1, $2, $3) RETURNING id"
	// Faz a transferência
	var transferID string
	err := r.QueryRow(ctx, query, accountOriginID, accountDestinationID, amount).Scan(&transferID)
	if err != nil {
		return entities.Transfer{}, err
	}

	transfer := entities.Transfer{
		ID:                   transferID,
		AccountOriginID:      accountOriginID,
		AccountDestinationID: accountDestinationID,
		Amount:               amount,
		CreatedAt:            time.Now(),
	}

	if err != nil {
		return entities.Transfer{}, err
	}

	return transfer, nil
}
