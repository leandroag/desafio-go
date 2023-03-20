package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/leandroag/desafio/app/domain/entities"
)

type accountRepository struct {
	Querier
}

func NewAccountRepository(querier Querier) *accountRepository {
	return &accountRepository{
		querier,
	}
}

func (r accountRepository) GetAccountByID(ctx context.Context, accountID int32) (entities.Account, error) {
	const query = "SELECT id, name, cpf, secret, balance, created_at FROM accounts WHERE id = $1"

	account := entities.Account{}
	err := r.QueryRow(ctx, query, accountID).
		Scan(&account.ID, &account.Name, &account.CPF, &account.Secret, &account.Balance, &account.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.Account{}, nil
		}
		return entities.Account{}, err
	}
	return account, nil
}

func (r accountRepository) GetAccountByCPF(ctx context.Context, CPF string) (entities.Account, error) {
	const query = "SELECT id, name, cpf, secret, balance, created_at FROM accounts WHERE cpf = $1"

	account := entities.Account{}
	err := r.QueryRow(ctx, query, CPF).
		Scan(&account.ID, &account.Name, &account.CPF, &account.Secret, &account.Balance, &account.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entities.Account{}, nil
		}
		return entities.Account{}, err
	}
	return account, nil
}

func (r accountRepository) GetAccountBalance(ctx context.Context, accountID int32) (float64, error) {
	const query = "SELECT balance FROM accounts WHERE id = $1"

	var balance float64
	err := r.QueryRow(ctx, query, accountID).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return balance, nil
}

func (r accountRepository) GetAllAccounts(ctx context.Context) ([]entities.Account, error) {
	const query = "SELECT id, name, cpf, secret, balance, created_at FROM accounts"
	rows, err := r.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []entities.Account{}
	for rows.Next() {
		account := entities.Account{}
		err = rows.Scan(&account.ID, &account.Name, &account.CPF, &account.Secret, &account.Balance, &account.CreatedAt)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r accountRepository) CreateAccount(ctx context.Context, account entities.Account) error {
	const query = "INSERT INTO accounts (name, cpf, secret, balance, created_at) VALUES ($1, $2, $3, $4, $5)"
	account.CreatedAt = time.Now()

	_, err := r.Exec(ctx, query, account.Name, account.CPF, account.Secret, account.Balance, account.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r accountRepository) UpdateAccount(ctx context.Context, account entities.Account) error {
	const query = `UPDATE accounts SET balance=$2 WHERE id=$1`
	_, err := r.Exec(ctx, query, account.ID, account.Balance)
	if err != nil {
		return err
	}
	return nil
}
