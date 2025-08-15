package repository

import (
	"wallet-api/internal/models"
	"wallet-api/utils"
)

type WalletRepositoryInterface interface {
	GetWalletByID(walletID string) (*models.Wallet, error)
	UpdateWalletBalance(walletID, operationType string, amount utils.Money) (*models.Wallet, error)
	CreateWallet(wallet *models.Wallet) error
}
