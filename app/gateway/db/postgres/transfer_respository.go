package postgres

import (
	"context"

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

func (r transferRepository) CreateTransfer(ctx context.Context, transfer entities.Transfer) error {
	const query = "INSERT INTO transfers(account_origin_id, account_destination_id, amount) VALUES($1, $2, $3) RETURNING id"

	err := r.QueryRow(ctx, query, transfer.AccountOriginID, transfer.AccountDestinationID, transfer.Amount).Scan(&transfer.ID)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}

func (r transferRepository) GetTransfersByAccountID(ctx context.Context, accountID int32) ([]entities.Transfer, error) {
	const query = "SELECT id, account_origin_id, account_destination_id, amount, created_at FROM transfers WHERE account_origin_id = $1"

	rows, err := r.Query(ctx, query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transfers := []entities.Transfer{}
	for rows.Next() {
		transfer := entities.Transfer{}
		err = rows.Scan(&transfer.ID, &transfer.AccountOriginID, &transfer.AccountDestinationID, &transfer.Amount, &transfer.CreatedAt)
		if err != nil {
			return nil, err
		}
		transfers = append(transfers, transfer)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return transfers, nil
}
