package handler

import (
	"encoding/json"
	stdErrors "errors"
	"fmt"
	"net/http"
	"wallet-api/internal/models"
	"wallet-api/internal/repository"
	"wallet-api/internal/service"
	"wallet-api/utils"

	"github.com/google/uuid"
)

// Временная структура для декодирования JSON с рубли
type walletOperationRequest struct {
	WalletID      uuid.UUID `json:"walletId"`
	OperationType string    `json:"operationType"`
	Amount        float64   `json:"amount"` // Рубли от пользователя
}

type WalletHandler struct {
	service service.WalletServiceInterface
}

func NewWalletHandler(service service.WalletServiceInterface) *WalletHandler {
	return &WalletHandler{service: service}
}

func (h *WalletHandler) HandleWalletOperation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var request walletOperationRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	// Конвертируем рубли в копейки
	amountStr := fmt.Sprintf("%.2f", request.Amount)
	money, err := utils.NewMoneyFromString(amountStr)
	if err != nil {
		http.Error(w, "Неверный формат суммы", http.StatusBadRequest)
		return
	}

	// Создаем операцию для сервиса с копейками
	operation := models.WalletOperation{
		WalletID:      request.WalletID,
		OperationType: request.OperationType,
		Amount:        money.Raw, // int64 в копейках
	}

	if err := h.validateWalletOperation(&operation); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	wallet, err := h.service.ProcessWalletOperation(&operation)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"walletId": wallet.ID,
		"balance":  wallet.Balance.String(),
	})
}

func (h *WalletHandler) validateWalletOperation(operation *models.WalletOperation) error {
	if operation.WalletID == uuid.Nil {
		return fmt.Errorf("ID кошелька обязателен")
	}

	if operation.OperationType == "" {
		return fmt.Errorf("Тип операции обязателен")
	}

	if !models.IsValidOperationType(operation.OperationType) {
		return fmt.Errorf("Неверный тип операции: %s", operation.OperationType)
	}

	if operation.Amount <= 0 {
		return fmt.Errorf("Сумма должна быть положительной")
	}

	return nil
}

func (h *WalletHandler) handleServiceError(w http.ResponseWriter, err error) {
	if stdErrors.Is(err, repository.ErrWalletNotFound) {
		http.Error(w, "Кошелек не найден", http.StatusNotFound)
		return
	}

	if stdErrors.Is(err, service.ErrInsufficientFunds) {
		http.Error(w, "Недостаточно средств", http.StatusBadRequest)
		return
	}

	http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
}

func (h *WalletHandler) HandleGetWallet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	walletID := r.URL.Path[len("/api/v1/wallets/"):]
	if walletID == "" {
		http.Error(w, "UUID кошелька не указан", http.StatusBadRequest)
		return
	}

	wallet, err := h.service.GetWallet(walletID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      wallet.ID,
		"balance": wallet.Balance.String(),
	})
}
