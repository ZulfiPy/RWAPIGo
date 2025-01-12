package main

import "time"

type Customer struct {
	FirstName      string
	LastName       string
	PersonalID     int64
	PhoneNumber    string
	Email          string
	RentedVehicles []string
	CreatedAt      time.Time
	LastEditedAt   *time.Time
}

type Customers []Customer
