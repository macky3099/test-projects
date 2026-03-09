package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func InitializeDatabase(databasePath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	createWalletsTableSQL := `
	CREATE TABLE IF NOT EXISTS wallets (
		wallet_id TEXT PRIMARY KEY,
		balance REAL NOT NULL DEFAULT 0.0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	createTransfersTableSQL := `
	CREATE TABLE IF NOT EXISTS transfers (
		transfer_id TEXT PRIMARY KEY,
		from_wallet_id TEXT NOT NULL,
		to_wallet_id TEXT NOT NULL,
		amount REAL NOT NULL,
		status TEXT NOT NULL DEFAULT 'PENDING',
		idempotency_key TEXT NOT NULL UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (from_wallet_id) REFERENCES wallets(wallet_id),
		FOREIGN KEY (to_wallet_id) REFERENCES wallets(wallet_id)
	);
	`

	createLedgerTableSQL := `
	CREATE TABLE IF NOT EXISTS ledger_entries (
		entry_id TEXT PRIMARY KEY,
		transfer_id TEXT NOT NULL,
		wallet_id TEXT NOT NULL,
		entry_type TEXT NOT NULL,
		amount REAL NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (transfer_id) REFERENCES transfers(transfer_id),
		FOREIGN KEY (wallet_id) REFERENCES wallets(wallet_id)
	);
	`

	createIdempotencyCacheTableSQL := `
	CREATE TABLE IF NOT EXISTS idempotency_cache (
		idempotency_key TEXT PRIMARY KEY,
		transfer_id TEXT NOT NULL,
		response_data TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (transfer_id) REFERENCES transfers(transfer_id)
	);
	`

	statements := []string{
		createWalletsTableSQL,
		createTransfersTableSQL,
		createLedgerTableSQL,
		createIdempotencyCacheTableSQL,
	}

	for _, sql := range statements {
		if _, err := db.Exec(sql); err != nil {
			return fmt.Errorf("failed to execute SQL: %w", err)
		}
	}

	return nil
}
