package repository

import (
	"database/sql"
	"fmt"
	"wallet-api/internal/models"
	"wallet-api/utils"
)

type WalletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) GetWalletByID(walletID string) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.QueryRow(
		`SELECT 
		id, 
		balance, 
		created_at, 
		updated_at 
		FROM wallets 
		WHERE id = $1`,
		walletID,
	).Scan(
		&wallet.ID,
		&wallet.Balance,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("get wallet by id: %w", ErrWalletNotFound)
		}
		return nil, fmt.Errorf("get wallet by id: %w", ErrDatabaseError)
	}

	return &wallet, nil
}

func (r *WalletRepository) UpdateWalletBalance(walletID string, operationType string, amount utils.Money) (*models.Wallet, error) {
	var wallet models.Wallet

	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", ErrDatabaseError)
	}
	defer tx.Rollback()

	updateQuery := `
		UPDATE wallets 
		SET balance = CASE 
			WHEN $2 = $3 THEN balance + $4
			WHEN $2 = $5 THEN balance - $4
		END,
		updated_at = NOW()
		WHERE id = $1
		RETURNING 
		id, 
		balance, 
		created_at, 
		updated_at
	`

	err = tx.QueryRow(
		updateQuery,
		walletID,
		operationType,
		models.OperationTypeDeposit,
		amount,
		models.OperationTypeWithdraw,
	).Scan(
		&wallet.ID,
		&wallet.Balance,
		&wallet.CreatedAt,
		&wallet.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			if operationType == models.OperationTypeDeposit {
				createQuery := `
					INSERT INTO wallets (id, balance, created_at, updated_at)
					VALUES ($1, $2, NOW(), NOW())
					RETURNING id, balance, created_at, updated_at
				`

				err = tx.QueryRow(
					createQuery,
					walletID,
					amount,
				).Scan(
					&wallet.ID,
					&wallet.Balance,
					&wallet.CreatedAt,
					&wallet.UpdatedAt,
				)

				if err != nil {
					return nil, fmt.Errorf("create wallet: %w", ErrDatabaseError)
				}
			} else {
				return nil, fmt.Errorf("update wallet balance: %w", ErrWalletNotFound)
			}
		} else {
			return nil, fmt.Errorf("update wallet balance: %w", ErrDatabaseError)
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", ErrDatabaseError)
	}

	return &wallet, nil
}

func (r *WalletRepository) CreateWallet(wallet *models.Wallet) error {
	query := `
		INSERT INTO wallets (
		id, 
		balance, 
		created_at, 
		updated_at
		)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.Exec(
		query,
		wallet.ID,
		wallet.Balance,
		wallet.CreatedAt,
		wallet.UpdatedAt)
	if err != nil {
		return fmt.Errorf("create wallet: %w", ErrDatabaseError)
	}

	return nil
}
