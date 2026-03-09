package service

import (
	"database/sql"
	"fmt"
	"time"
	"wallet-transfer/internal/models"
	"wallet-transfer/internal/repository"
)

type TransferService struct {
	db                    *sql.DB
	walletRepo            *repository.WalletRepository
	transferRepo          *repository.TransferRepository
	ledgerRepo            *repository.LedgerRepository
	idempotencyRepo       *repository.IdempotencyRepository
}

func NewTransferService(db *sql.DB) *TransferService {
	return &TransferService{
		db:              db,
		walletRepo:      repository.NewWalletRepository(db),
		transferRepo:    repository.NewTransferRepository(db),
		ledgerRepo:      repository.NewLedgerRepository(db),
		idempotencyRepo: repository.NewIdempotencyRepository(db),
	}
}

func (ts *TransferService) ProcessTransfer(request *models.TransferRequest) (*models.Transfer, error) {
	if request.FromWalletID == request.ToWalletID {
		return nil, fmt.Errorf("cannot transfer to the same wallet")
	}

	if request.Amount <= 0 {
		return nil, fmt.Errorf("transfer amount must be positive")
	}

	transferID, alreadyProcessed, err := ts.idempotencyRepo.CheckIfProcessed(request.IdempotencyKey)
	if err != nil {
		return nil, err
	}

	if alreadyProcessed {
		transfer, err := ts.transferRepo.GetTransferByID(transferID)
		if err != nil {
			return nil, err
		}
		return transfer, nil
	}

	tx, err := ts.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	fromBalance, err := ts.walletRepo.GetWalletBalanceWithLock(tx, request.FromWalletID)
	if err != nil {
		return nil, err
	}

	if fromBalance < request.Amount {
		return nil, fmt.Errorf("insufficient balance in wallet %s", request.FromWalletID)
	}

	_, err = ts.walletRepo.GetWalletBalanceWithLock(tx, request.ToWalletID)
	if err != nil {
		return nil, err
	}

	newTransferID := generateTransferID()
	transfer := &models.Transfer{
		TransferID:     newTransferID,
		FromWalletID:   request.FromWalletID,
		ToWalletID:     request.ToWalletID,
		Amount:         request.Amount,
		Status:         "PROCESSED",
		IdempotencyKey: request.IdempotencyKey,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = ts.transferRepo.CreateTransfer(transfer)
	if err != nil {
		return nil, err
	}

	err = ts.walletRepo.UpdateWalletBalance(tx, request.FromWalletID, fromBalance-request.Amount)
	if err != nil {
		return nil, err
	}

	toBalance, err := ts.walletRepo.GetWalletBalanceWithLock(tx, request.ToWalletID)
	if err != nil {
		return nil, err
	}

	err = ts.walletRepo.UpdateWalletBalance(tx, request.ToWalletID, toBalance+request.Amount)
	if err != nil {
		return nil, err
	}

	debitEntry := &models.LedgerEntry{
		EntryID:    generateEntryID(),
		TransferID: newTransferID,
		WalletID:   request.FromWalletID,
		EntryType:  "DEBIT",
		Amount:     request.Amount,
		CreatedAt:  time.Now(),
	}

	creditEntry := &models.LedgerEntry{
		EntryID:    generateEntryID(),
		TransferID: newTransferID,
		WalletID:   request.ToWalletID,
		EntryType:  "CREDIT",
		Amount:     request.Amount,
		CreatedAt:  time.Now(),
	}

	err = ts.ledgerRepo.AddLedgerEntry(debitEntry)
	if err != nil {
		return nil, err
	}

	err = ts.ledgerRepo.AddLedgerEntry(creditEntry)
	if err != nil {
		return nil, err
	}

	err = ts.idempotencyRepo.SaveIdempotencyRecord(request.IdempotencyKey, newTransferID, "success")
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return transfer, nil
}

func (ts *TransferService) GetWalletBalance(walletID string) (float64, error) {
	balance, err := ts.walletRepo.GetWalletBalance(walletID)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func (ts *TransferService) GetWalletLedger(walletID string) ([]*models.LedgerEntry, error) {
	entries, err := ts.ledgerRepo.GetLedgerEntriesByWallet(walletID)
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func (ts *TransferService) EnsureWalletExists(walletID string) error {
	return ts.EnsureWalletExistsWithBalance(walletID, 0)
}

func (ts *TransferService) EnsureWalletExistsWithBalance(walletID string, initialBalance float64) error {
	_, err := ts.walletRepo.GetWalletBalance(walletID)
	if err != nil {
		err := ts.walletRepo.CreateWallet(walletID, initialBalance)
		if err != nil {
			return fmt.Errorf("failed to create wallet: %w", err)
		}
	}
	return nil
}

func generateTransferID() string {
	return fmt.Sprintf("TXN_%d", time.Now().UnixNano())
}

func generateEntryID() string {
	return fmt.Sprintf("ENTRY_%d", time.Now().UnixNano())
}
