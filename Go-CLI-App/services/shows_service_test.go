package services

import (
	"booking_system/respository"
	"errors"
	"testing"
)

type mockRepo struct{}

func (m mockRepo) AvailableSeats(key int) ([]respository.Seat, bool) {
	return []respository.Seat{
		{ID: "B1", Type: "B", Price: 280.0, Reserved: false},
	}, true
}

func (m mockRepo) BookSeats(seats map[string][]string, shownumber int, seatnumbers []string) ([]respository.Seat, error) {
	if shownumber != 1 {
		return nil, errors.New("invalid show")
	}
	return []respository.Seat{
		{ID: "B1", Type: "B", Price: 280.0, Reserved: true},
	}, nil
}

func (m mockRepo) SaveCurrentSale(revenue float64, allTaxes map[string]float64) {}
func (m mockRepo) GetTotalSalesSummary() respository.Sales {
	return respository.Sales{Revenue: 1000, ServiceTax: 140}
}

func TestGetAllSeats(t *testing.T) {
	service := NewShowService(mockRepo{})

	seats, err := service.GetAllSeats(1)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(seats["B"]) != 1 {
		t.Errorf("expected 1 seat in B, got %d", len(seats["B"]))
	}
}

func TestBookSeats(t *testing.T) {
	service := NewShowService(mockRepo{})

	seatMap := map[string][]string{
		"B": {"B1"},
	}
	service.BookSeats(seatMap, 1, []string{"B1"})
	// No panic means success for this test (side-effect based)
}

func TestGetRevenueDetails(t *testing.T) {
	service := NewShowService(mockRepo{})
	revenue := service.GetRevenueDetails()

	if revenue.Revenue != 1000 {
		t.Errorf("expected 1000 revenue, got %.2f", revenue.Revenue)
	}
}
