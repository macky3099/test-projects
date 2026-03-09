package models

import "time"

type Wallet struct {
	WalletID string
	Balance  float64
}

type Transfer struct {
	TransferID    string
	FromWalletID  string
	ToWalletID    string
	Amount        float64
	Status        string
	IdempotencyKey string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type LedgerEntry struct {
	EntryID      string
	TransferID   string
	WalletID     string
	EntryType    string
	Amount       float64
	CreatedAt    time.Time
}

type TransferRequest struct {
	FromWalletID   string  `json:"fromWalletID"`
	ToWalletID     string  `json:"toWalletID"`
	Amount         float64 `json:"amount"`
	IdempotencyKey string  `json:"idempotencyKey"`
}

type TransferResponse struct {
	TransferID string `json:"transferID"`
	Status     string `json:"status"`
	Message    string `json:"message"`
}

type BalanceResponse struct {
	WalletID string  `json:"walletID"`
	Balance  float64 `json:"balance"`
}

type LedgerResponse struct {
	EntryID    string    `json:"entryID"`
	TransferID string    `json:"transferID"`
	WalletID   string    `json:"walletID"`
	EntryType  string    `json:"entryType"`
	Amount     float64   `json:"amount"`
	CreatedAt  time.Time `json:"createdAt"`
}
