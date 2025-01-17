package customer

import (
	"errors"
	"fmt"
	"net/mail"
	"time"
	"unicode"

	"github.com/ZulfiPy/RWAPIGo/internal/storage"
)

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

type CustomerStorage struct {
	store *storage.Storage[Customers]
}

func NewCustomerStorage(fileName string) *CustomerStorage {
	return &CustomerStorage{
		store: storage.NewStorage[Customers](fileName),
	}
}

func NewCustomerStorage1(fileName string) *storage.Storage[Customers] {
	return storage.NewStorage[Customers](fileName)
}

func (cs *CustomerStorage) Load(data *Customers) error {
	return cs.store.Load(data)
}

func (cs *CustomerStorage) GetStorage() *storage.Storage[Customers] {
	return cs.store
}

func intLength(number int64) int {
	if number == 0 {
		return 1
	}

	length := 0

	for number != 0 {
		number /= 10
		length++
	}

	return length
}

func (cs CustomerStorage) ValidateInput(input Customer) error {
	if input.FirstName == "" || len(input.FirstName) < 3 {
		return errors.New("invalid input: first name cannot be empty or shorter than 3 characters")
	}

	if input.LastName == "" || len(input.LastName) < 3 {
		return errors.New("invalid input: last name cannot be empty or shorter than 3 characters")
	}

	if input.PhoneNumber == "" || len(input.PhoneNumber) < 7 {
		return errors.New("invalid input: phone number cannot be empty or shorter than 7 numbers")
	}

	for _, char := range input.PhoneNumber {
		isChar := unicode.IsLetter(char)

		if isChar {
			return errors.New("invalid input: phone number cannot consist letters")
		}
	}

	if input.Email == "" || len(input.Email) < 7 {
		return errors.New("invalid input: email cannot be empty or shorter than 7 characters")
	}

	_, err := mail.ParseAddress(input.Email)

	if err != nil {
		return fmt.Errorf("invalid %v", err)
	}

	personalIDLen := intLength(input.PersonalID)

	if personalIDLen != 11 {
		return errors.New("invalid input: personal id of the customer must be exactly 11 digits")
	}

	return nil
}
