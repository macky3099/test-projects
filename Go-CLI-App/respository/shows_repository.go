package respository

import (
	"errors"
)

type Repository interface {
	AvailableSeats(key int) ([]Seat, bool)
	BookSeats(seats map[string][]string, shownumber int, seatnumbers []string) ([]Seat, error)
	SaveCurrentSale(Revenue float64, allTaxes map[string]float64)
	GetTotalSalesSummary() Sales
}

type baseRepo struct {
}

func NewRepository() Repository {

	repo := &baseRepo{}

	return repo
}

func (r *baseRepo) AvailableSeats(showID int) ([]Seat, bool) {
	seats, exists := AvailableSeats[showID]
	if !exists {
		return nil, false
	}
	return seats, true
}

func (r *baseRepo) BookSeats(seats map[string][]string, shownumber int, seatnumbers []string) ([]Seat, error) {
	// Check if show exists
	showSeats, exists := AvailableSeats[shownumber]
	if !exists {
		return nil, errors.New("Unable to book any show at the moment, please try later.")
	}

	var bookedSeats []Seat

	// Iterate over requested seat numbers
	for _, requestedSeat := range seatnumbers {
		for i := range showSeats {
			if showSeats[i].ID == requestedSeat && !showSeats[i].Reserved {
				// Mark seat as reserved
				showSeats[i].Reserved = true
				bookedSeats = append(bookedSeats, showSeats[i])
			}
		}
	}

	// Update the global AvailableSeats map
	AvailableSeats[shownumber] = showSeats

	return bookedSeats, nil
}

func (r *baseRepo) SaveCurrentSale(revenue float64, allTaxes map[string]float64) {
	var currentSales Sales
	currentSales.Revenue = revenue

	// Map the taxes to struct fields based on the keys in your tax map
	for taxType, amount := range allTaxes {
		switch taxType {
		case "Service Tax":
			currentSales.ServiceTax = amount
		case "Swachh Bharat Cess":
			currentSales.SwachBharatCess = amount
		case "Krishi Kalyan Cess":
			currentSales.KrishiKalyanCess = amount
		}
	}

	// Append this sale to the global slice
	TotalSales = append(TotalSales, currentSales)
}

func (r *baseRepo) GetTotalSalesSummary() Sales {
	var summary Sales

	for _, sale := range TotalSales {
		summary.Revenue += sale.Revenue
		summary.ServiceTax += sale.ServiceTax
		summary.SwachBharatCess += sale.SwachBharatCess
		summary.KrishiKalyanCess += sale.KrishiKalyanCess
	}

	return summary
}
