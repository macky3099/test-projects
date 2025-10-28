package respository

import (
	"testing"
)

func TestBookSeats(t *testing.T) {
	// Setup: Prepare mock available seats
	AvailableSeats = map[int][]Seat{
		1: {
			{ID: "A1", Type: "A", Price: 100.0, Reserved: false},
			{ID: "A2", Type: "A", Price: 100.0, Reserved: false},
		},
	}

	repo := NewRepository()

	seatsToBook := []string{"A1"}
	bookedSeats, err := repo.BookSeats(nil, 1, seatsToBook)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(bookedSeats) != 1 {
		t.Errorf("expected 1 booked seat, got %d", len(bookedSeats))
	}

	if !AvailableSeats[1][0].Reserved {
		t.Errorf("expected A1 to be reserved")
	}
}

func TestGetTotalSalesSummary(t *testing.T) {
	// Setup: Reset TotalSales
	TotalSales = []Sales{
		{Revenue: 100, ServiceTax: 10, SwachBharatCess: 1, KrishiKalyanCess: 1},
		{Revenue: 200, ServiceTax: 20, SwachBharatCess: 2, KrishiKalyanCess: 2},
	}

	repo := NewRepository()
	summary := repo.GetTotalSalesSummary()

	if summary.Revenue != 300 {
		t.Errorf("expected revenue 300, got %.2f", summary.Revenue)
	}
	if summary.ServiceTax != 30 {
		t.Errorf("expected service tax 30, got %.2f", summary.ServiceTax)
	}
}
