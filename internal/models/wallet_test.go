package models

import (
	"testing"
	"time"
	"wallet-api/utils"

	"github.com/google/uuid"
)

func TestIsValidOperationType(t *testing.T) {
	tests := []struct {
		name     string
		opType   string
		expected bool
	}{
		{"valid deposit", OperationTypeDeposit, true},
		{"valid withdraw", OperationTypeWithdraw, true},
		{"invalid operation", "TRANSFER", false},
		{"empty string", "", false},
		{"case sensitive", "deposit", false},
		{"case sensitive withdraw", "withdraw", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidOperationType(tt.opType)
			if result != tt.expected {
				t.Errorf("IsValidOperationType(%s) = %v, want %v", tt.opType, result, tt.expected)
			}
		})
	}
}

func TestWalletOperation_Validation(t *testing.T) {
	validWalletID := uuid.New()

	tests := []struct {
		name    string
		op      WalletOperation
		isValid bool
	}{
		{
			name: "valid deposit operation",
			op: WalletOperation{
				WalletID:      validWalletID,
				OperationType: OperationTypeDeposit,
				Amount:        1000,
			},
			isValid: true,
		},
		{
			name: "valid withdraw operation",
			op: WalletOperation{
				WalletID:      validWalletID,
				OperationType: OperationTypeWithdraw,
				Amount:        500,
			},
			isValid: true,
		},
		{
			name: "invalid operation type",
			op: WalletOperation{
				WalletID:      validWalletID,
				OperationType: "INVALID",
				Amount:        1000,
			},
			isValid: false,
		},
		{
			name: "zero amount",
			op: WalletOperation{
				WalletID:      validWalletID,
				OperationType: OperationTypeDeposit,
				Amount:        0,
			},
			isValid: true,
		},
		{
			name: "negative amount",
			op: WalletOperation{
				WalletID:      validWalletID,
				OperationType: OperationTypeDeposit,
				Amount:        -100,
			},
			isValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := IsValidOperationType(tt.op.OperationType)
			if isValid != tt.isValid {
				t.Errorf("WalletOperation validation = %v, want %v", isValid, tt.isValid)
			}
		})
	}
}

func TestWallet_Initialization(t *testing.T) {
	walletID := uuid.New()
	balance := utils.Money{Raw: 1000}

	wallet := &Wallet{
		ID:        walletID,
		Balance:   balance,
		CreatedAt: time.Now(),
	}

	if wallet.ID != walletID {
		t.Errorf("Wallet ID = %v, want %v", wallet.ID, walletID)
	}

	if wallet.Balance.Raw != balance.Raw {
		t.Errorf("Wallet Balance = %v, want %v", wallet.Balance.Raw, balance.Raw)
	}

	if wallet.CreatedAt.IsZero() {
		t.Error("Wallet CreatedAt should not be zero")
	}
}
