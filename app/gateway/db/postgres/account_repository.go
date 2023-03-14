package postgres

import (
	"database/sql"
	"time"

	"github.com/leandroag/desafio/app/domain/entities"
)

type accountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *accountRepository {
	return &accountRepository{
		db: db,
	}
}

func (repository accountRepository) GetAccountByID(accountID int64) (*entities.Account, error) {
	account := &entities.Account{}
	err := repository.db.QueryRow("SELECT id, name, cpf, secret, balance, created_at FROM accounts WHERE id = $1", accountID).
		Scan(&account.ID, &account.Name, &account.CPF, &account.Secret, &account.Balance, &account.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return account, nil
}

func (repository accountRepository) GetAccountBalance(accountID int64) (float64, error) {
	var balance float64
	err := repository.db.QueryRow("SELECT balance FROM accounts WHERE id = $1", accountID).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return balance, nil
}

func (repository accountRepository) GetAllAccounts() ([]*entities.Account, error) {
	rows, err := repository.db.Query("SELECT id, name, cpf, secret, balance, created_at FROM accounts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	accounts := []*entities.Account{}
	for rows.Next() {
		account := &entities.Account{}
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

func (repo accountRepository) CreateAccount(account entities.Account) error {
	account.CreatedAt = time.Now()

	_, err := repo.db.Exec("INSERT INTO accounts (name, cpf, secret, balance, created_at) VALUES ($1, $2, $3, $4, $5)",
		account.Name, account.CPF, account.Secret, account.Balance, account.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
