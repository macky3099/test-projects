package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"wallet-transfer/internal/database"
	"wallet-transfer/internal/handler"
	"wallet-transfer/internal/models"
	"wallet-transfer/internal/service"
)

func setupTestDB() *service.TransferService {
	db, _ := database.InitializeDatabase(":memory:")
	return service.NewTransferService(db)
}

func TestHandlePostTransferSuccess(t *testing.T) {
	transferService := setupTestDB()
	transferHandler := handler.NewTransferHandler(transferService)

	transferService.EnsureWalletExistsWithBalance("wallet1", 1000)
	transferService.EnsureWalletExistsWithBalance("wallet2", 0)

	requestBody := models.TransferRequest{
		FromWalletID:   "wallet1",
		ToWalletID:     "wallet2",
		Amount:         100,
		IdempotencyKey: "key1",
	}

	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest("POST", "/transfers", bytes.NewReader(body))
	w := httptest.NewRecorder()

	transferHandler.HandlePostTransfer(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response models.TransferResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Status != "PROCESSED" {
		t.Errorf("expected status PROCESSED, got %s", response.Status)
	}
}

func TestHandlePostTransferInvalidRequest(t *testing.T) {
	transferService := setupTestDB()
	transferHandler := handler.NewTransferHandler(transferService)

	req := httptest.NewRequest("POST", "/transfers", bytes.NewReader([]byte("invalid")))
	w := httptest.NewRecorder()

	transferHandler.HandlePostTransfer(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHandleGetBalance(t *testing.T) {
	transferService := setupTestDB()
	transferHandler := handler.NewTransferHandler(transferService)

	transferService.EnsureWalletExists("wallet1")

	req := httptest.NewRequest("GET", "/balance?walletID=wallet1", nil)
	w := httptest.NewRecorder()

	transferHandler.HandleGetBalance(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response models.BalanceResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.WalletID != "wallet1" {
		t.Errorf("expected walletID wallet1, got %s", response.WalletID)
	}
}

func TestHandleGetBalanceMissingWalletID(t *testing.T) {
	transferService := setupTestDB()
	transferHandler := handler.NewTransferHandler(transferService)

	req := httptest.NewRequest("GET", "/balance", nil)
	w := httptest.NewRecorder()

	transferHandler.HandleGetBalance(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestHandleGetLedger(t *testing.T) {
	transferService := setupTestDB()
	transferHandler := handler.NewTransferHandler(transferService)

	transferService.EnsureWalletExists("wallet1")

	req := httptest.NewRequest("GET", "/ledger?walletID=wallet1", nil)
	w := httptest.NewRecorder()

	transferHandler.HandleGetLedger(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
