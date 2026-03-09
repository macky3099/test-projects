package repository

import (
	"database/sql"
	"fmt"
	"time"
	"wallet-transfer/internal/models"
)

type WalletRepository struct {
	db *sql.DB
}

type TransferRepository struct {
	db *sql.DB
}

type LedgerRepository struct {
	db *sql.DB
}

type IdempotencyRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func NewTransferRepository(db *sql.DB) *TransferRepository {
	return &TransferRepository{db: db}
}

func NewLedgerRepository(db *sql.DB) *LedgerRepository {
	return &LedgerRepository{db: db}
}

func NewIdempotencyRepository(db *sql.DB) *IdempotencyRepository {
	return &IdempotencyRepository{db: db}
}

func (wr *WalletRepository) CreateWallet(walletID string, initialBalance float64) error {
	query := `INSERT INTO wallets (wallet_id, balance) VALUES (?, ?)`
	_, err := wr.db.Exec(query, walletID, initialBalance)
	if err != nil {
		return fmt.Errorf("failed to create wallet: %w", err)
	}
	return nil
}

func (wr *WalletRepository) GetWalletBalance(walletID string) (float64, error) {
	query := `SELECT balance FROM wallets WHERE wallet_id = ?`
	var balance float64
	err := wr.db.QueryRow(query, walletID).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("wallet not found: %s", walletID)
		}
		return 0, fmt.Errorf("failed to get wallet balance: %w", err)
	}
	return balance, nil
}

func (wr *WalletRepository) GetWalletBalanceWithLock(tx *sql.Tx, walletID string) (float64, error) {
	query := `SELECT balance FROM wallets WHERE wallet_id = ? FOR UPDATE`
	var balance float64
	err := tx.QueryRow(query, walletID).Scan(&balance)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("wallet not found: %s", walletID)
		}
		return 0, fmt.Errorf("failed to get wallet balance: %w", err)
	}
	return balance, nil
}

func (wr *WalletRepository) UpdateWalletBalance(tx *sql.Tx, walletID string, newBalance float64) error {
	query := `UPDATE wallets SET balance = ?, updated_at = CURRENT_TIMESTAMP WHERE wallet_id = ?`
	_, err := tx.Exec(query, newBalance, walletID)
	if err != nil {
		return fmt.Errorf("failed to update wallet balance: %w", err)
	}
	return nil
}

func (tr *TransferRepository) CreateTransfer(transfer *models.Transfer) error {
	query := `
	INSERT INTO transfers (transfer_id, from_wallet_id, to_wallet_id, amount, status, idempotency_key, created_at, updated_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := tr.db.Exec(query,
		transfer.TransferID,
		transfer.FromWalletID,
		transfer.ToWalletID,
		transfer.Amount,
		transfer.Status,
		transfer.IdempotencyKey,
		transfer.CreatedAt,
		transfer.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create transfer: %w", err)
	}
	return nil
}

func (tr *TransferRepository) GetTransferByID(transferID string) (*models.Transfer, error) {
	query := `
	SELECT transfer_id, from_wallet_id, to_wallet_id, amount, status, idempotency_key, created_at, updated_at
	FROM transfers WHERE transfer_id = ?
	`
	var transfer models.Transfer
	err := tr.db.QueryRow(query, transferID).Scan(
		&transfer.TransferID,
		&transfer.FromWalletID,
		&transfer.ToWalletID,
		&transfer.Amount,
		&transfer.Status,
		&transfer.IdempotencyKey,
		&transfer.CreatedAt,
		&transfer.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transfer not found: %s", transferID)
		}
		return nil, fmt.Errorf("failed to get transfer: %w", err)
	}
	return &transfer, nil
}

func (tr *TransferRepository) UpdateTransferStatus(tx *sql.Tx, transferID string, status string) error {
	query := `UPDATE transfers SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE transfer_id = ?`
	_, err := tx.Exec(query, status, transferID)
	if err != nil {
		return fmt.Errorf("failed to update transfer status: %w", err)
	}
	return nil
}

func (lr *LedgerRepository) AddLedgerEntry(entry *models.LedgerEntry) error {
	query := `
	INSERT INTO ledger_entries (entry_id, transfer_id, wallet_id, entry_type, amount, created_at)
	VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := lr.db.Exec(query,
		entry.EntryID,
		entry.TransferID,
		entry.WalletID,
		entry.EntryType,
		entry.Amount,
		entry.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to add ledger entry: %w", err)
	}
	return nil
}

func (lr *LedgerRepository) GetLedgerEntriesByWallet(walletID string) ([]*models.LedgerEntry, error) {
	query := `
	SELECT entry_id, transfer_id, wallet_id, entry_type, amount, created_at
	FROM ledger_entries WHERE wallet_id = ? ORDER BY created_at DESC
	`
	rows, err := lr.db.Query(query, walletID)
	if err != nil {
		return nil, fmt.Errorf("failed to get ledger entries: %w", err)
	}
	defer rows.Close()

	var entries []*models.LedgerEntry
	for rows.Next() {
		var entry models.LedgerEntry
		err := rows.Scan(
			&entry.EntryID,
			&entry.TransferID,
			&entry.WalletID,
			&entry.EntryType,
			&entry.Amount,
			&entry.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan ledger entry: %w", err)
		}
		entries = append(entries, &entry)
	}

	return entries, nil
}

func (ir *IdempotencyRepository) CheckIfProcessed(idempotencyKey string) (string, bool, error) {
	query := `SELECT transfer_id FROM idempotency_cache WHERE idempotency_key = ?`
	var transferID string
	err := ir.db.QueryRow(query, idempotencyKey).Scan(&transferID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", false, nil
		}
		return "", false, fmt.Errorf("failed to check idempotency: %w", err)
	}
	return transferID, true, nil
}

func (ir *IdempotencyRepository) SaveIdempotencyRecord(idempotencyKey string, transferID string, responseData string) error {
	query := `
	INSERT INTO idempotency_cache (idempotency_key, transfer_id, response_data, created_at)
	VALUES (?, ?, ?, ?)
	`
	_, err := ir.db.Exec(query, idempotencyKey, transferID, responseData, time.Now())
	if err != nil {
		return fmt.Errorf("failed to save idempotency record: %w", err)
	}
	return nil
}
