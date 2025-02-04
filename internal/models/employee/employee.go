package employee

import (
	"errors"
	"fmt"

	"github.com/ZulfiPy/RWAPIGo/internal/storage"
	"github.com/ZulfiPy/RWAPIGo/internal/utils"
)

type Employee struct {
	FirstName   string
	LastName    string
	PersonalID  int64
	DateOfBirth string `json:"DateOfBirth"`
	Email       string
	PhoneNumber string
	Address     string
}

type Employees []Employee

type EmployeeStorage struct {
	storage *storage.Storage[Employees]
}

func (es *EmployeeStorage) validateInput(input Employee) error {
	if input.FirstName == "" || len(input.FirstName) < 3 {
		return errors.New("invalid input: first name cannot be empty or shorter than 3 characters")
	}

	if input.LastName == "" || len(input.LastName) < 3 {
		return errors.New("invalid input: last name cannot be empty or shorter than 3 characters")
	}

	personalIDLen := utils.IntLength(input.PersonalID)

	if personalIDLen != 11 {
		return errors.New("invalid input: personal id of the employee must be exactly 11 digits")
	}

	if !utils.IsValidDateFormat(input.DateOfBirth) {
		return errors.New("invalid input: wrong date format")
	}

	emailErr := utils.IsValidEmail(input.Email)

	if emailErr != nil {
		return emailErr
	}

	if input.PhoneNumber == "" || len(input.PhoneNumber) < 7 {
		return errors.New("invalid input: phone number cannot be empty or shorter than 7 numbers")
	}

	if input.Address == "" || len(input.Address) < 5 {
		return errors.New("invalid input: living address cannot be empty or shorter than 5 symbols")
	}

	return nil
}

func (es *EmployeeStorage) duplicatedEmployee(employees Employees, personalID int64) bool {
	for _, employee := range employees {
		if employee.PersonalID == personalID {
			return true
		}
	}

	return false
}

func NewEmployeeStorage(fileName string) *EmployeeStorage {
	return &EmployeeStorage{
		storage: storage.NewStorage[Employees](fileName),
	}
}

func (es *EmployeeStorage) GetStorage() *storage.Storage[Employees] {
	return es.storage
}

func (es *EmployeeStorage) GetEmployees() (Employees, error) {
	employees := Employees{}

	if err := es.storage.Load(&employees); err != nil {
		return Employees{}, err
	}

	return employees, nil
}

func (es *EmployeeStorage) AddEmployee(input Employee) (Employee, error) {
	employees := Employees{}
	if err := es.storage.Load(&employees); err != nil {
		return Employee{}, err
	}

	if err := es.validateInput(input); err != nil {
		return Employee{}, err
	}

	if es.duplicatedEmployee(employees, input.PersonalID) {
		return Employee{}, fmt.Errorf("employee with personal ID %d already exists", input.PersonalID)
	}

	employees = append(employees, input)
	if err := es.storage.Save(employees); err != nil {
		return Employee{}, err
	}

	return input, nil
}
