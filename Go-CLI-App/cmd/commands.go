package cmd

import (
	"booking_system/handlers"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func StartBooking(handler *handlers.ShowsBaseHandler) {
	for {
		showNumber := GetShow()
		seats, shownumber := handler.ShowAllSeats(showNumber)

		seatNumbers := GetSeatNumbers()
		handler.BookSeats(seats, shownumber, seatNumbers)

		var choice string
		fmt.Print("\n\nWould you like to continue (Yes/No): ")
		fmt.Scanln(&choice)

		if strings.ToLower(choice) != "yes" {
			handler.ShowRevenue()
			break
		}
	}
}

func GetShow() int {
	reader := bufio.NewReader(os.Stdin)
	var number int

	for {
		fmt.Print("\n\nEnter Show no: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)
		number, err = strconv.Atoi(input)
		if err != nil {
			fmt.Println("\nInvalid number, please try again.")
			continue
		}

		// Valid number entered
		break
	}

	return number
}

func GetSeatNumbers() []string {
	reader := bufio.NewReader(os.Stdin)
	var seats []string

	for {
		fmt.Print("\nEnter seats: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			fmt.Println("Input cannot be empty, please try again.")
			continue
		}

		// Split by comma and trim spaces
		rawSeats := strings.Split(input, ",")
		for _, s := range rawSeats {
			trimmed := strings.TrimSpace(s)
			if trimmed != "" {
				seats = append(seats, trimmed)
			}
		}

		if len(seats) == 0 {
			fmt.Println("No valid seats entered, try again.")
			continue
		}

		break
	}
	return seats
}
