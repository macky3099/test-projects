package handlers

import (
	"booking_system/respository"
	"testing"
)

// Mock service to test handler behavior
type mockShowService struct {
	getSeatsFunc      func(int) (map[string][]string, error)
	bookSeatsFunc     func(map[string][]string, int, []string)
	getRevenueDetails func() respository.Sales
}

func (m mockShowService) GetAllSeats(shownumber int) (map[string][]string, error) {
	return m.getSeatsFunc(shownumber)
}

func (m mockShowService) BookSeats(seats map[string][]string, shownumber int, seatnumbers []string) {
	m.bookSeatsFunc(seats, shownumber, seatnumbers)
}

func (m mockShowService) GetRevenueDetails() respository.Sales {
	return m.getRevenueDetails()
}

func TestShowAllSeats(t *testing.T) {
	mockService := mockShowService{
		getSeatsFunc: func(show int) (map[string][]string, error) {
			return map[string][]string{
				"A": {"A1", "A2"},
				"B": {"B1"},
				"C": {},
			}, nil
		},
	}
	handler := NewShowsHandler(mockService)

	seats, show := handler.ShowAllSeats(1)
	if show != 1 {
		t.Errorf("expected show number 1, got %d", show)
	}
	if len(seats["A"]) != 2 {
		t.Errorf("expected 2 seats in A, got %d", len(seats["A"]))
	}
}

func TestBookSeatsUnavailable(t *testing.T) {
	mockService := mockShowService{
		bookSeatsFunc: func(seats map[string][]string, shownumber int, seatnumbers []string) {},
	}
	handler := NewShowsHandler(mockService)

	seats := map[string][]string{
		"A": {"A1"},
		"B": {"B1"},
		"C": {"C1"},
	}

	ok, err := handler.BookSeats(seats, 1, []string{"Z9"}) // invalid seat
	if ok {
		t.Error("expected booking to fail, but it succeeded")
	}
	if err == nil {
		t.Error("expected an error, got nil")
	}
}

func TestBookSeatsSuccess(t *testing.T) {
	mockService := mockShowService{
		bookSeatsFunc: func(seats map[string][]string, shownumber int, seatnumbers []string) {},
	}
	handler := NewShowsHandler(mockService)

	seats := map[string][]string{
		"A": {"A1"},
	}

	ok, err := handler.BookSeats(seats, 1, []string{"A1"})
	if !ok {
		t.Error("expected booking to succeed")
	}
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}
