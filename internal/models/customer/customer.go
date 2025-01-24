package customer

import (
	"errors"
	"fmt"
	"net/mail"
	"time"
	"unicode"

	"github.com/ZulfiPy/RWAPIGo/internal/models/vehicle"
	"github.com/ZulfiPy/RWAPIGo/internal/storage"
)

type Customer struct {
	FirstName      string
	LastName       string
	PersonalID     int64
	PhoneNumber    string
	Email          string
	RentedVehicles []vehicle.Vehicle
	CreatedAt      time.Time
	LastEditedAt   *time.Time
}

type Customers []Customer

type CustomerStorage struct {
	storage *storage.Storage[Customers]
}

func NewCustomerStorage(fileName string) *CustomerStorage {
	return &CustomerStorage{
		storage: storage.NewStorage[Customers](fileName),
	}
}

func (cs *CustomerStorage) GetStorage() *storage.Storage[Customers] {
	return cs.storage
}

func IntLength(number int64) int {
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

func validateIndex(idx, customersLength int) error {
	if idx < 0 || idx >= customersLength {
		return fmt.Errorf("error:index %d is out of range", idx)
	}

	return nil
}

func (cs *CustomerStorage) validateInput(input Customer) error {
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

	personalIDLen := IntLength(input.PersonalID)

	if personalIDLen != 11 {
		return errors.New("invalid input: personal id of the customer must be exactly 11 digits")
	}

	return nil
}

func (cs *CustomerStorage) findCustomerByPersonalID(personalID int64) int {
	customers := Customers{}
	cs.storage.Load(&customers)

	for idx, customer := range customers {
		if customer.PersonalID == personalID {
			return idx
		}
	}

	return -1
}

func (cs *CustomerStorage) AddCustomer(input Customer) error {
	if err := cs.validateInput(input); err != nil {
		return err
	}

	if idx := cs.findCustomerByPersonalID(input.PersonalID); idx != -1 {
		return fmt.Errorf("customer with personalID %d is found in the storage, duplicated customers not allowed", input.PersonalID)
	}

	customers := &Customers{}
	cs.storage.Load(customers)

	newCustomer := Customer{
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		PersonalID:     input.PersonalID,
		PhoneNumber:    input.PhoneNumber,
		Email:          input.Email,
		RentedVehicles: []vehicle.Vehicle{},
		CreatedAt:      time.Now(),
	}

	*customers = append(*customers, newCustomer)

	if err := cs.storage.Save(*customers); err != nil {
		return err
	}

	return nil
}

func (cs *CustomerStorage) DeleteCustomer(personalID int64) error {
	customers := Customers{}
	cs.storage.Load(&customers)

	idx := cs.findCustomerByPersonalID(personalID)

	if idx == -1 {
		return fmt.Errorf("customer with personalID %d not found", personalID)
	}

	if err := validateIndex(idx, len(customers)); err != nil {
		return err
	}

	customers = append(customers[:idx], customers[idx+1:]...)

	if err := cs.storage.Save(customers); err != nil {
		return err
	}

	return nil
}

func (cs *CustomerStorage) EditCustomer(firstName, lastName, email, phoneNumber string, personalID int64) error {
	idx := cs.findCustomerByPersonalID(personalID)

	if idx == -1 {
		return fmt.Errorf("customer with personalID %d not found in the storage", personalID)
	}

	customers := Customers{}
	cs.storage.Load(&customers)

	customerToEdit := &customers[idx]

	if len(firstName) != 0 {
		customerToEdit.FirstName = firstName
	}

	if len(lastName) != 0 {
		customerToEdit.LastName = lastName
	}

	if len(email) != 0 {
		customerToEdit.Email = email
	}

	if len(phoneNumber) != 0 {
		customerToEdit.PhoneNumber = phoneNumber
	}

	lastEdited := time.Now()

	customerToEdit.LastEditedAt = &lastEdited

	if err := cs.storage.Save(customers); err != nil {
		return err
	}

	return nil
}

func (cs *CustomerStorage) GetCustomers() (Customers, error) {
	customers := Customers{}

	if err := cs.storage.Load(&customers); err != nil {
		return nil, err
	}

	return customers, nil
}

func (cs *CustomerStorage) AddVehicle(vehicle vehicle.Vehicle, personalID int64) (Customer, error) {
	customers := Customers{}
	if err := cs.storage.Load(&customers); err != nil {
		return Customer{}, err
	}

	idx := cs.findCustomerByPersonalID(personalID)
	if idx == -1 {
		return Customer{}, fmt.Errorf("customer with personalID %d not found", personalID)
	}

	customers[idx].RentedVehicles = append(customers[idx].RentedVehicles, vehicle)

	if err := cs.storage.Save(customers); err != nil {
		return Customer{}, err
	}
	return customers[idx], nil
}

func (cs *CustomerStorage) DeleteVehicle(plateNumber string, personalID int64) error {
	customers := Customers{}
	if err := cs.storage.Load(&customers); err != nil {
		return err
	}

	customerIdx := cs.findCustomerByPersonalID(personalID)
	if customerIdx == -1 {
		return fmt.Errorf("customer with personalID %d not found", personalID)
	}

	for idx, vehicle := range customers[customerIdx].RentedVehicles {
		if vehicle.PlateNumber == plateNumber {
			customers[customerIdx].RentedVehicles = append(customers[customerIdx].RentedVehicles[:idx], customers[customerIdx].RentedVehicles[idx+1:]...)
			break
		}
	}

	if err := cs.storage.Save(customers); err != nil {
		return err
	}

	return nil
}
