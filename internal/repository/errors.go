package repository

import "errors"

var (
	ErrWalletNotFound = errors.New("wallet not found in repository")
	ErrDatabaseError  = errors.New("database error")
)
