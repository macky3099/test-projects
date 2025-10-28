package handlers

import (
	"booking_system/services"
	"errors"
	"fmt"
	"strings"
)

type ShowsBaseHandler struct {
	service services.ShowsService
}

func NewShowsHandler(services services.ShowsService) *ShowsBaseHandler {
	return &ShowsBaseHandler{
		service: services,
	}
}

func (b *ShowsBaseHandler) ShowAllSeats(shownumber int) (map[string][]string, int) {
	groupedSeats, err := b.service.GetAllSeats(shownumber)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Available Seats:")
		fmt.Println(strings.Join(groupedSeats["A"], " "))
		fmt.Println(strings.Join(groupedSeats["B"], " "))
		fmt.Println(strings.Join(groupedSeats["C"], " "))
	}
	return groupedSeats, shownumber
}

func (b *ShowsBaseHandler) BookSeats(seats map[string][]string, shownumber int, seatnumbers []string) (bool, error) {
	// Iterate over requested seat numbers
	for _, requestedSeat := range seatnumbers {
		found := false

		// Check each row (A, B, C...)
		for _, rowSeats := range seats {
			for _, availableSeat := range rowSeats {
				if requestedSeat == availableSeat {
					found = true
					break
				}
			}
			if found {
				b.service.BookSeats(seats, shownumber, seatnumbers)
				return true, nil
			}
		}

		if !found {
			fmt.Printf("Print: %s NOT available, Please select different seats\n", requestedSeat)
			return false, errors.New(fmt.Sprintf("Seat %s is NOT available for show %d.\n", requestedSeat, shownumber))
		}
	}
	return false, fmt.Errorf("Some error occured")
}

func (b *ShowsBaseHandler) ShowRevenue() {
	totalRevenue := b.service.GetRevenueDetails()
	fmt.Printf("\nTotal Sales:\n")
	fmt.Printf("Revenue: Rs.%.2f\n", totalRevenue.Revenue)
	fmt.Printf("Service Tax: Rs.%.2f\n", totalRevenue.ServiceTax)
	fmt.Printf("Swachh Bharat Cess: Rs.%.2f\n", totalRevenue.SwachBharatCess)
	fmt.Printf("Krishi Kalyan Cess: Rs.%.2f\n", totalRevenue.KrishiKalyanCess)

	fmt.Print("\n")
}
