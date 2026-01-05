package main

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func main() {
	var conferenceName = "Go Conference"
	const conferenceTicket = 50
	var remainingTickets uint = 50
	bookings := []string{}

	fmt.Println("Hello World Welcome to our conference", conferenceName)
	fmt.Printf("We have total of %v Ticket where remaining is %v \n", conferenceTicket, remainingTickets)

	var userName string
	var userTickets uint
	var userEmail string
	// ask for user namer.

	type UserDetail struct {
		username string
		ticket   int
		Password string
	}
	userDetail := UserDetail{
		username: "Ritik",
		ticket:   2,
		Password: "hirit",
	}
	fmt.Printf("User Detail:\n", userDetail)

	for {

		if remainingTickets == 0 {
			fmt.Println("Our conference is booked now")
			break
		}

		fmt.Printf("Enter user name:\n")
		fmt.Scan(&userName)
		bookings = append(bookings, userName)

		fmt.Printf("Enter user ticket:\n")
		fmt.Scan(&userTickets)

		fmt.Printf("Enter user email:\n")
		fmt.Scan(&userEmail)

		isValidName := len(userName) >= 2
		isValidEmail := strings.Contains(userEmail, "@")

		if isValidName && !isValidEmail {
			fmt.Print("Invalid data")
			continue
		}

		remainingTickets = remainingTickets - userTickets
		fmt.Printf("Thank you %v for booking %v  tickets. You will receive a confirmation email at %v - %v\n", userName, userTickets, userEmail, uuid.New())

		fmt.Println("Remaining Ticket is ", remainingTickets, bookings)
	}

}
