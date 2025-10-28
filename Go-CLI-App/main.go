package main

import (
	"booking_system/cmd"
	"booking_system/handlers"
	"booking_system/respository"
	"booking_system/services"
)

func main() {
	//init db
	//pass to repo when db is used

	repo := respository.NewRepository()
	service := services.NewShowService(repo)
	handler := handlers.NewShowsHandler(service)

	cmd.StartBooking(handler)
}
