package models

import (
	"database/sql"
	"time"
	"wallet-api/utils"

	"github.com/google/uuid"
)

type Wallet struct {
	ID        uuid.UUID    `db:"id" json:"id"`
	Balance   utils.Money  `db:"balance" json:"balance"`
	CreatedAt time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at" json:"updated_at"`
}

type WalletOperation struct {
	WalletID      uuid.UUID `json:"walletId"`
	OperationType string    `json:"operationType"`
	Amount        int64     `json:"amount"`
}

const (
	OperationTypeDeposit  = "DEPOSIT"
	OperationTypeWithdraw = "WITHDRAW"
)

func IsValidOperationType(opType string) bool {
	return opType == OperationTypeDeposit || opType == OperationTypeWithdraw
}
