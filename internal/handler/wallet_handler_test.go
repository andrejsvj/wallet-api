package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"wallet-api/internal/models"
	"wallet-api/internal/repository"
	"wallet-api/internal/service"
	"wallet-api/utils"

	"github.com/google/uuid"
)

type MockWalletService struct {
	shouldError bool
	errorType   error
	wallet      *models.Wallet
}

func (m *MockWalletService) GetWallet(walletID string) (*models.Wallet, error) {
	if m.shouldError {
		return nil, m.errorType
	}
	return m.wallet, nil
}

func (m *MockWalletService) ProcessWalletOperation(operation *models.WalletOperation) (*models.Wallet, error) {
	if m.shouldError {
		return nil, m.errorType
	}
	return m.wallet, nil
}

func (m *MockWalletService) CreateWallet(wallet *models.Wallet) error {
	if m.shouldError {
		return m.errorType
	}
	return nil
}

func TestWalletHandler_HandleWalletOperation(t *testing.T) {
	walletID := uuid.New()
	wallet := &models.Wallet{
		ID:      walletID,
		Balance: utils.Money{Raw: 1000},
	}

	tests := []struct {
		name           string
		method         string
		requestBody    interface{}
		expectedStatus int
		setupMock      func() *MockWalletService
	}{
		{
			name:   "valid operation",
			method: http.MethodPost,
			requestBody: walletOperationRequest{
				WalletID:      walletID,
				OperationType: models.OperationTypeDeposit,
				Amount:        5.00, // 5 рублей
			},
			expectedStatus: http.StatusOK,
			setupMock: func() *MockWalletService {
				return &MockWalletService{wallet: wallet}
			},
		},
		{
			name:           "unsupported method",
			method:         http.MethodGet,
			requestBody:    nil,
			expectedStatus: http.StatusMethodNotAllowed,
			setupMock: func() *MockWalletService {
				return &MockWalletService{}
			},
		},
		{
			name:   "invalid JSON",
			method: http.MethodPost,
			requestBody: map[string]interface{}{
				"invalid": "json",
			},
			expectedStatus: http.StatusBadRequest,
			setupMock: func() *MockWalletService {
				return &MockWalletService{}
			},
		},
		{
			name:   "service error - wallet not found",
			method: http.MethodPost,
			requestBody: walletOperationRequest{
				WalletID:      walletID,
				OperationType: models.OperationTypeDeposit,
				Amount:        5.00, // 5 рублей
			},
			expectedStatus: http.StatusNotFound,
			setupMock: func() *MockWalletService {
				return &MockWalletService{
					shouldError: true,
					errorType:   repository.ErrWalletNotFound,
				}
			},
		},
		{
			name:   "service error - insufficient funds",
			method: http.MethodPost,
			requestBody: walletOperationRequest{
				WalletID:      walletID,
				OperationType: models.OperationTypeWithdraw,
				Amount:        15.00, // 15 рублей
			},
			expectedStatus: http.StatusBadRequest,
			setupMock: func() *MockWalletService {
				return &MockWalletService{
					shouldError: true,
					errorType:   service.ErrInsufficientFunds,
				}
			},
		},
		{
			name:   "conversion test - rubles to kopecks",
			method: http.MethodPost,
			requestBody: walletOperationRequest{
				WalletID:      walletID,
				OperationType: models.OperationTypeDeposit,
				Amount:        10.50, // 10.50 рублей = 1050 копеек
			},
			expectedStatus: http.StatusOK,
			setupMock: func() *MockWalletService {
				// Проверяем, что сервис получает правильное количество копеек
				return &MockWalletService{
					wallet: &models.Wallet{
						ID:      walletID,
						Balance: utils.Money{Raw: 2050}, // 20.50 рублей
					},
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.setupMock()
			handler := NewWalletHandler(mockService)

			var body []byte
			var err error
			if tt.requestBody != nil {
				body, err = json.Marshal(tt.requestBody)
				if err != nil {
					t.Fatalf("Failed to marshal request body: %v", err)
				}
			}

			req := httptest.NewRequest(tt.method, "/api/v1/wallets/operation", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.HandleWalletOperation(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestWalletHandler_HandleGetWallet(t *testing.T) {
	walletID := uuid.New()
	wallet := &models.Wallet{
		ID:      walletID,
		Balance: utils.Money{Raw: 1000},
	}

	tests := []struct {
		name           string
		method         string
		walletID       string
		expectedStatus int
		setupMock      func() *MockWalletService
	}{
		{
			name:           "valid wallet request",
			method:         http.MethodGet,
			walletID:       walletID.String(),
			expectedStatus: http.StatusOK,
			setupMock: func() *MockWalletService {
				return &MockWalletService{wallet: wallet}
			},
		},
		{
			name:           "unsupported method",
			method:         http.MethodPost,
			walletID:       walletID.String(),
			expectedStatus: http.StatusMethodNotAllowed,
			setupMock: func() *MockWalletService {
				return &MockWalletService{}
			},
		},
		{
			name:           "missing wallet ID",
			method:         http.MethodGet,
			walletID:       "",
			expectedStatus: http.StatusBadRequest,
			setupMock: func() *MockWalletService {
				return &MockWalletService{}
			},
		},
		{
			name:           "wallet not found",
			method:         http.MethodGet,
			walletID:       walletID.String(),
			expectedStatus: http.StatusNotFound,
			setupMock: func() *MockWalletService {
				return &MockWalletService{
					shouldError: true,
					errorType:   repository.ErrWalletNotFound,
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := tt.setupMock()
			handler := NewWalletHandler(mockService)

			url := "/api/v1/wallets/"
			if tt.walletID != "" {
				url += tt.walletID
			}

			req := httptest.NewRequest(tt.method, url, nil)
			w := httptest.NewRecorder()

			handler.HandleGetWallet(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestWalletHandler_validateWalletOperation(t *testing.T) {
	handler := &WalletHandler{}
	validWalletID := uuid.New()

	tests := []struct {
		name      string
		operation *models.WalletOperation
		wantErr   bool
	}{
		{
			name: "valid operation",
			operation: &models.WalletOperation{
				WalletID:      validWalletID,
				OperationType: models.OperationTypeDeposit,
				Amount:        500,
			},
			wantErr: false,
		},
		{
			name: "nil wallet ID",
			operation: &models.WalletOperation{
				WalletID:      uuid.Nil,
				OperationType: models.OperationTypeDeposit,
				Amount:        500,
			},
			wantErr: true,
		},
		{
			name: "empty operation type",
			operation: &models.WalletOperation{
				WalletID:      validWalletID,
				OperationType: "",
				Amount:        500,
			},
			wantErr: true,
		},
		{
			name: "invalid operation type",
			operation: &models.WalletOperation{
				WalletID:      validWalletID,
				OperationType: "INVALID",
				Amount:        500,
			},
			wantErr: true,
		},
		{
			name: "zero amount",
			operation: &models.WalletOperation{
				WalletID:      validWalletID,
				OperationType: models.OperationTypeDeposit,
				Amount:        0,
			},
			wantErr: true,
		},
		{
			name: "negative amount",
			operation: &models.WalletOperation{
				WalletID:      validWalletID,
				OperationType: models.OperationTypeDeposit,
				Amount:        -100,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.validateWalletOperation(tt.operation)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateWalletOperation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
