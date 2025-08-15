package service

import (
	"wallet-api/internal/models"
)

type WalletServiceInterface interface {
	GetWallet(walletID string) (*models.Wallet, error)
	ProcessWalletOperation(operation *models.WalletOperation) (*models.Wallet, error)
	CreateWallet(wallet *models.Wallet) error
}
