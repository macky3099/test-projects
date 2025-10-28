package services

import (
	"booking_system/respository"
	"fmt"
)

type showsBaseService struct {
	repo respository.Repository
}

type ShowsService interface {
	GetAllSeats(shownumber int) (map[string][]string, error)
	BookSeats(seats map[string][]string, shownumber int, seatnumbers []string)
	GetRevenueDetails() respository.Sales
}

func NewShowService(repo respository.Repository) ShowsService {
	return &showsBaseService{
		repo: repo,
	}
}

func (s showsBaseService) GetAllSeats(shownumber int) (map[string][]string, error) {
	seats, exists := s.repo.AvailableSeats(shownumber)
	if !exists {
		return nil, fmt.Errorf("show %d not found", shownumber)
	}

	grouped := map[string][]string{"A": {}, "B": {}, "C": {}}
	for _, seat := range seats {
		if seat.Price == 0 && seat.Reserved {
			// Skip seats which we don't want to include like C1 in Show 1
			continue
		}
		if !seat.Reserved {
			grouped[seat.Type] = append(grouped[seat.Type], seat.ID)
		}
	}
	return grouped, nil
}

func (s showsBaseService) BookSeats(seats map[string][]string, shownumber int, seatnumbers []string) {
	bookedSeats, _ := s.repo.BookSeats(seats, shownumber, seatnumbers)
	fmt.Printf("Print: Sucessfully Booked - Show %d\n\n", shownumber)
	s.showAndSaveSale(bookedSeats)
}

func (s showsBaseService) showAndSaveSale(bookedSeats []respository.Seat) {

	var SubTotal float64
	var TotalBill float64
	taxesApplied := make(map[string]float64)

	for _, seats := range bookedSeats {
		SubTotal = SubTotal + seats.Price
	}

	TotalBill = SubTotal
	fmt.Printf("Subtotal: Rs.%.2f\n", SubTotal)

	for taxType, rate := range respository.Taxes {
		total := SubTotal * float64(rate) / 100
		taxesApplied[taxType] = total
		TotalBill = TotalBill + total
		fmt.Printf("%s @%.2f%%: Rs.%.2f\n", taxType, rate, total)
	}
	fmt.Printf("Total: Rs.%.2f", TotalBill)
	s.repo.SaveCurrentSale(SubTotal, taxesApplied)
}

func (s showsBaseService) GetRevenueDetails() respository.Sales {
	return s.repo.GetTotalSalesSummary()
}
