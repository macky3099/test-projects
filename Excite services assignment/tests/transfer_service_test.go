package main

import (
	"testing"
	"wallet-transfer/internal/database"
	"wallet-transfer/internal/models"
	"wallet-transfer/internal/service"
)

func setupServiceForTest() *service.TransferService {
	db, _ := database.InitializeDatabase(":memory:")
	return service.NewTransferService(db)
}

func TestProcessTransferSuccess(t *testing.T) {
	transferService := setupServiceForTest()

	transferService.EnsureWalletExists("wallet1")
	transferService.EnsureWalletExists("wallet2")

	request := &models.TransferRequest{
		FromWalletID:   "wallet1",
		ToWalletID:     "wallet2",
		Amount:         100,
		IdempotencyKey: "key1",
	}

	transfer, err := transferService.ProcessTransfer(request)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if transfer.Status != "PROCESSED" {
		t.Errorf("expected status PROCESSED, got %s", transfer.Status)
	}
}

func TestProcessTransferInsufficientBalance(t *testing.T) {
	transferService := setupServiceForTest()

	transferService.EnsureWalletExists("wallet1")
	transferService.EnsureWalletExists("wallet2")

	request := &models.TransferRequest{
		FromWalletID:   "wallet1",
		ToWalletID:     "wallet2",
		Amount:         1000,
		IdempotencyKey: "key1",
	}

	_, err := transferService.ProcessTransfer(request)
	if err == nil {
		t.Errorf("expected error for insufficient balance, got nil")
	}
}

func TestProcessTransferIdempotency(t *testing.T) {
	transferService := setupServiceForTest()

	transferService.EnsureWalletExists("wallet1")
	transferService.EnsureWalletExists("wallet2")

	request := &models.TransferRequest{
		FromWalletID:   "wallet1",
		ToWalletID:     "wallet2",
		Amount:         100,
		IdempotencyKey: "key1",
	}

	transfer1, _ := transferService.ProcessTransfer(request)
	transfer2, _ := transferService.ProcessTransfer(request)

	if transfer1.TransferID != transfer2.TransferID {
		t.Errorf("expected same transfer ID for idempotent request, got different IDs")
	}
}

func TestProcessTransferSameWallet(t *testing.T) {
	transferService := setupServiceForTest()

	transferService.EnsureWalletExists("wallet1")

	request := &models.TransferRequest{
		FromWalletID:   "wallet1",
		ToWalletID:     "wallet1",
		Amount:         100,
		IdempotencyKey: "key1",
	}

	_, err := transferService.ProcessTransfer(request)
	if err == nil {
		t.Errorf("expected error for same wallet transfer, got nil")
	}
}

func TestProcessTransferNegativeAmount(t *testing.T) {
	transferService := setupServiceForTest()

	transferService.EnsureWalletExists("wallet1")
	transferService.EnsureWalletExists("wallet2")

	request := &models.TransferRequest{
		FromWalletID:   "wallet1",
		ToWalletID:     "wallet2",
		Amount:         -100,
		IdempotencyKey: "key1",
	}

	_, err := transferService.ProcessTransfer(request)
	if err == nil {
		t.Errorf("expected error for negative amount, got nil")
	}
}

func TestGetWalletBalance(t *testing.T) {
	transferService := setupServiceForTest()

	transferService.EnsureWalletExists("wallet1")

	balance, err := transferService.GetWalletBalance("wallet1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if balance != 0 {
		t.Errorf("expected balance 0, got %f", balance)
	}
}

func TestGetWalletLedger(t *testing.T) {
	transferService := setupServiceForTest()

	transferService.EnsureWalletExists("wallet1")
	transferService.EnsureWalletExists("wallet2")

	request := &models.TransferRequest{
		FromWalletID:   "wallet1",
		ToWalletID:     "wallet2",
		Amount:         100,
		IdempotencyKey: "key1",
	}

	transferService.ProcessTransfer(request)

	ledger, err := transferService.GetWalletLedger("wallet1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(ledger) == 0 {
		t.Errorf("expected ledger entries, got empty")
	}

	if ledger[0].EntryType != "DEBIT" {
		t.Errorf("expected entry type DEBIT, got %s", ledger[0].EntryType)
	}
}
