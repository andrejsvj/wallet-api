package service

import (
	"errors"
	"testing"
	"time"
	"wallet-api/internal/models"
	"wallet-api/internal/repository"
	"wallet-api/utils"
	"wallet-api/utils/logger"

	"github.com/google/uuid"
)

func init() {
	logger.Init()
}

type MockWalletRepository struct {
	wallets     map[string]*models.Wallet
	shouldError bool
	errorType   error
}

func NewMockWalletRepository() *MockWalletRepository {
	return &MockWalletRepository{
		wallets: make(map[string]*models.Wallet),
	}
}

func (m *MockWalletRepository) GetWalletByID(walletID string) (*models.Wallet, error) {
	if m.shouldError {
		return nil, m.errorType
	}

	wallet, exists := m.wallets[walletID]
	if !exists {
		return nil, repository.ErrWalletNotFound
	}
	return wallet, nil
}

func (m *MockWalletRepository) UpdateWalletBalance(walletID, operationType string, amount utils.Money) (*models.Wallet, error) {
	if m.shouldError {
		return nil, m.errorType
	}

	wallet, exists := m.wallets[walletID]
	if !exists {
		return nil, repository.ErrWalletNotFound
	}

	switch operationType {
	case models.OperationTypeDeposit:
		wallet.Balance = wallet.Balance.Add(amount)
	case models.OperationTypeWithdraw:
		wallet.Balance = wallet.Balance.Sub(amount)
	}

	wallet.UpdatedAt.Time = time.Now()
	wallet.UpdatedAt.Valid = true
	return wallet, nil
}

func (m *MockWalletRepository) CreateWallet(wallet *models.Wallet) error {
	if m.shouldError {
		return m.errorType
	}

	m.wallets[wallet.ID.String()] = wallet
	return nil
}

func TestWalletService_GetWallet(t *testing.T) {
	mockRepo := NewMockWalletRepository()
	service := NewWalletService(mockRepo)

	walletID := uuid.New()
	wallet := &models.Wallet{
		ID:        walletID,
		Balance:   utils.Money{Raw: 1000},
		CreatedAt: time.Now(),
	}
	mockRepo.wallets[walletID.String()] = wallet

	tests := []struct {
		name     string
		walletID string
		want     *models.Wallet
		wantErr  error
	}{
		{
			name:     "existing wallet",
			walletID: walletID.String(),
			want:     wallet,
			wantErr:  nil,
		},
		{
			name:     "non-existing wallet",
			walletID: uuid.New().String(),
			want:     nil,
			wantErr:  repository.ErrWalletNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetWallet(tt.walletID)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("WalletService.GetWallet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr == nil && got.ID != tt.want.ID {
				t.Errorf("WalletService.GetWallet() = %v, want %v", got.ID, tt.want.ID)
			}
		})
	}
}

func TestWalletService_ProcessWalletOperation(t *testing.T) {
	mockRepo := NewMockWalletRepository()
	service := NewWalletService(mockRepo)

	walletID := uuid.New()
	wallet := &models.Wallet{
		ID:        walletID,
		Balance:   utils.Money{Raw: 1000},
		CreatedAt: time.Now(),
	}
	mockRepo.wallets[walletID.String()] = wallet

	tests := []struct {
		name      string
		operation *models.WalletOperation
		wantErr   error
		setupMock func()
	}{
		{
			name: "valid deposit operation",
			operation: &models.WalletOperation{
				WalletID:      walletID,
				OperationType: models.OperationTypeDeposit,
				Amount:        500,
			},
			wantErr: nil,
		},
		{
			name: "valid withdraw operation with sufficient funds",
			operation: &models.WalletOperation{
				WalletID:      walletID,
				OperationType: models.OperationTypeWithdraw,
				Amount:        300,
			},
			wantErr: nil,
		},
		{
			name: "withdraw operation with insufficient funds",
			operation: &models.WalletOperation{
				WalletID:      walletID,
				OperationType: models.OperationTypeWithdraw,
				Amount:        1500,
			},
			wantErr: ErrInsufficientFunds,
		},

		{
			name: "non-existing wallet",
			operation: &models.WalletOperation{
				WalletID:      uuid.New(),
				OperationType: models.OperationTypeDeposit,
				Amount:        100,
			},
			wantErr: repository.ErrWalletNotFound,
		},
		{
			name: "database error during withdraw check",
			operation: &models.WalletOperation{
				WalletID:      walletID,
				OperationType: models.OperationTypeWithdraw,
				Amount:        100,
			},
			wantErr: repository.ErrDatabaseError,
			setupMock: func() {
				mockRepo.shouldError = true
				mockRepo.errorType = repository.ErrDatabaseError
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			} else {
				mockRepo.shouldError = false
			}

			_, err := service.ProcessWalletOperation(tt.operation)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("WalletService.ProcessWalletOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWalletService_ProcessWalletOperation_BalanceChanges(t *testing.T) {
	tests := []struct {
		name            string
		operation       *models.WalletOperation
		expectedBalance int64
	}{
		{
			name: "deposit increases balance",
			operation: &models.WalletOperation{
				WalletID:      uuid.New(),
				OperationType: models.OperationTypeDeposit,
				Amount:        500,
			},
			expectedBalance: 1500,
		},
		{
			name: "withdraw decreases balance",
			operation: &models.WalletOperation{
				WalletID:      uuid.New(),
				OperationType: models.OperationTypeWithdraw,
				Amount:        300,
			},
			expectedBalance: 700,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := NewMockWalletRepository()
			service := NewWalletService(mockRepo)

			initialBalance := utils.Money{Raw: 1000}
			wallet := &models.Wallet{
				ID:        tt.operation.WalletID,
				Balance:   initialBalance,
				CreatedAt: time.Now(),
			}
			mockRepo.wallets[tt.operation.WalletID.String()] = wallet

			result, err := service.ProcessWalletOperation(tt.operation)
			if err != nil {
				t.Errorf("WalletService.ProcessWalletOperation() unexpected error = %v", err)
				return
			}

			if result.Balance.Raw != tt.expectedBalance {
				t.Errorf("WalletService.ProcessWalletOperation() balance = %v, want %v",
					result.Balance.Raw, tt.expectedBalance)
			}
		})
	}
}

func TestWalletService_CreateWallet(t *testing.T) {
	mockRepo := NewMockWalletRepository()
	service := NewWalletService(mockRepo)

	wallet := &models.Wallet{
		ID:        uuid.New(),
		Balance:   utils.Money{Raw: 0},
		CreatedAt: time.Now(),
	}

	tests := []struct {
		name      string
		wallet    *models.Wallet
		wantErr   error
		setupMock func()
	}{
		{
			name:    "successful wallet creation",
			wallet:  wallet,
			wantErr: nil,
		},
		{
			name:    "database error during creation",
			wallet:  wallet,
			wantErr: repository.ErrDatabaseError,
			setupMock: func() {
				mockRepo.shouldError = true
				mockRepo.errorType = repository.ErrDatabaseError
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupMock != nil {
				tt.setupMock()
			} else {
				mockRepo.shouldError = false
			}

			err := service.CreateWallet(tt.wallet)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("WalletService.CreateWallet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
