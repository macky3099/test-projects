package handler

import (
	"encoding/json"
	"net/http"
	"wallet-transfer/internal/models"
	"wallet-transfer/internal/service"
)

type TransferHandler struct {
	transferService *service.TransferService
}

func NewTransferHandler(transferService *service.TransferService) *TransferHandler {
	return &TransferHandler{
		transferService: transferService,
	}
}

func (th *TransferHandler) HandlePostTransfer(w http.ResponseWriter, r *http.Request) {
	var request models.TransferRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
		return
	}

	if request.FromWalletID == "" || request.ToWalletID == "" || request.Amount <= 0 || request.IdempotencyKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "missing or invalid fields"})
		return
	}

	err = th.transferService.EnsureWalletExists(request.FromWalletID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	err = th.transferService.EnsureWalletExists(request.ToWalletID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	transfer, err := th.transferService.ProcessTransfer(&request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := models.TransferResponse{
		TransferID: transfer.TransferID,
		Status:     transfer.Status,
		Message:    "transfer successful",
	}
	json.NewEncoder(w).Encode(response)
}

func (th *TransferHandler) HandleGetBalance(w http.ResponseWriter, r *http.Request) {
	walletID := r.URL.Query().Get("walletID")
	if walletID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "walletID is required"})
		return
	}

	balance, err := th.transferService.GetWalletBalance(walletID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := models.BalanceResponse{
		WalletID: walletID,
		Balance:  balance,
	}
	json.NewEncoder(w).Encode(response)
}

func (th *TransferHandler) HandleGetLedger(w http.ResponseWriter, r *http.Request) {
	walletID := r.URL.Query().Get("walletID")
	if walletID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "walletID is required"})
		return
	}

	entries, err := th.transferService.GetWalletLedger(walletID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var response []models.LedgerResponse
	for _, entry := range entries {
		response = append(response, models.LedgerResponse{
			EntryID:    entry.EntryID,
			TransferID: entry.TransferID,
			WalletID:   entry.WalletID,
			EntryType:  entry.EntryType,
			Amount:     entry.Amount,
			CreatedAt:  entry.CreatedAt,
		})
	}
	json.NewEncoder(w).Encode(response)
}

func SetupRoutes(handler *TransferHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /transfers", handler.HandlePostTransfer)
	mux.HandleFunc("GET /balance", handler.HandleGetBalance)
	mux.HandleFunc("GET /ledger", handler.HandleGetLedger)
	return mux
}
