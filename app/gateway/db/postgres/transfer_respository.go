package postgres

import (
	"database/sql"
	"time"

	"github.com/leandroag/desafio/app/domain/entities"
)

type transferRepository struct {
	db *sql.DB
}

func NewTransferRepository(db *sql.DB) *transferRepository {
	return &transferRepository{
		db: db,
	}
}

func (repository transferRepository) CreateTransfer(accountOriginID string, accountDestinationID string, amount float64) (entities.Transfer, error) {
	tx, err := repository.db.Begin()
	if err != nil {
		return entities.Transfer{}, err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Faz a transferência
	var transferID string
	err = tx.QueryRow("INSERT INTO transfers(account_origin_id, account_destination_id, amount) VALUES($1, $2, $3) RETURNING id", accountOriginID, accountDestinationID, amount).Scan(&transferID)
	if err != nil {
		tx.Rollback()
		return entities.Transfer{}, err
	}

	transfer := entities.Transfer{
		ID:                   transferID,
		AccountOriginID:      accountOriginID,
		AccountDestinationID: accountDestinationID,
		Amount:               amount,
		CreatedAt:            time.Now(),
	}

	err = tx.Commit()
	if err != nil {
		return entities.Transfer{}, err
	}

	return transfer, nil
}
