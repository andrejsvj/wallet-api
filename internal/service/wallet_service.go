package service

import (
	stdErrors "errors"
	"fmt"
	"wallet-api/internal/models"
	"wallet-api/internal/repository"
	"wallet-api/utils"
	"wallet-api/utils/logger"
)

type WalletService struct {
	repo repository.WalletRepositoryInterface
}

func NewWalletService(repo repository.WalletRepositoryInterface) *WalletService {
	return &WalletService{repo: repo}
}

func (s *WalletService) GetWallet(walletID string) (*models.Wallet, error) {
	wallet, err := s.repo.GetWalletByID(walletID)
	if err != nil {
		if stdErrors.Is(err, repository.ErrWalletNotFound) {
			return nil, fmt.Errorf("get wallet: %w", repository.ErrWalletNotFound)
		}
		return nil, fmt.Errorf("get wallet: %w", repository.ErrDatabaseError)
	}
	return wallet, nil
}

func (s *WalletService) ProcessWalletOperation(operation *models.WalletOperation) (*models.Wallet, error) {
	if operation.OperationType == models.OperationTypeWithdraw {
		existingWallet, err := s.repo.GetWalletByID(operation.WalletID.String())
		if err != nil {
			if stdErrors.Is(err, repository.ErrWalletNotFound) {
				return nil, fmt.Errorf("process operation: %w", repository.ErrWalletNotFound)
			}
			return nil, fmt.Errorf("process operation: %w", repository.ErrDatabaseError)
		}

		withdrawAmount := utils.Money{Raw: operation.Amount}
		if existingWallet.Balance.Raw < withdrawAmount.Raw {
			logger.GlobalLogger.Warning("Insufficient funds detected for wallet %s", existingWallet.ID)
			return nil, fmt.Errorf("process operation: %w", ErrInsufficientFunds)
		}
	}

	amount := utils.Money{Raw: operation.Amount}
	wallet, err := s.repo.UpdateWalletBalance(operation.WalletID.String(), operation.OperationType, amount)
	if err != nil {
		if stdErrors.Is(err, repository.ErrWalletNotFound) {
			return nil, fmt.Errorf("process operation: %w", repository.ErrWalletNotFound)
		}
		return nil, fmt.Errorf("process operation: %w", repository.ErrDatabaseError)
	}

	return wallet, nil
}

func (s *WalletService) CreateWallet(wallet *models.Wallet) error {
	err := s.repo.CreateWallet(wallet)
	if err != nil {
		return fmt.Errorf("create wallet: %w", repository.ErrDatabaseError)
	}

	return nil
}
