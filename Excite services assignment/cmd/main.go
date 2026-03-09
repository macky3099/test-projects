package main

import (
	"fmt"
	"log"
	"net/http"
	"wallet-transfer/internal/database"
	"wallet-transfer/internal/handler"
	"wallet-transfer/internal/service"
)

func main() {
	db, err := database.InitializeDatabase("wallet_transfer.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	transferService := service.NewTransferService(db)
	transferHandler := handler.NewTransferHandler(transferService)
	router := handler.SetupRoutes(transferHandler)

	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(port, router))
}
