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

	if input.Address == "" || len(input.Address) < 7 {
		return errors.New("invalid input: living address cannot be empty or shorter than 7 symbols")
	}

	return nil
}

func (es *EmployeeStorage) validateEditData(email, phoneNumber, address string) error {
	emailErr := utils.IsValidEmail(email)

	if emailErr != nil {
		return emailErr
	}

	if phoneNumber == "" || len(phoneNumber) < 7 {
		return errors.New("invalid input: phone number cannot be empty or shorter than 7 numnbers")
	}

	if address == "" || len(address) < 7 {
		return errors.New("invalid input: living address cannot be empty or shorter than 7 symbols")
	}

	return nil
}

func (es *EmployeeStorage) employeePersists(employees Employees, personalID int64) (int, error) {
	for idx, employee := range employees {
		if employee.PersonalID == personalID {
			return idx, nil
		}
	}

	return -1, fmt.Errorf("employee with personalID %d not found", personalID)
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

	if _, err := es.employeePersists(employees, input.PersonalID); err != nil {
		return Employee{}, fmt.Errorf("employee with personal ID %d already exists", input.PersonalID)
	}

	employees = append(employees, input)
	if err := es.storage.Save(employees); err != nil {
		return Employee{}, err
	}

	return input, nil
}

func (es *EmployeeStorage) DeleteEmployee(personalID int64) error {
	employees := Employees{}
	if err := es.storage.Load(&employees); err != nil {
		return err
	}

	if _, err := es.employeePersists(employees, personalID); err != nil {
		return err
	}

	for idx, employee := range employees {
		if employee.PersonalID == personalID {
			employees = append(employees[:idx], employees[idx+1:]...)
			break
		}
	}

	if err := es.storage.Save(employees); err != nil {
		return err
	}

	return nil
}

func (es *EmployeeStorage) EditEmployeeContacts(email, phoneNumber, address string, personalID int64) (Employee, error) {
	employees := Employees{}

	if err := es.storage.Load(&employees); err != nil {
		return Employee{}, err
	}

	idx, err := es.employeePersists(employees, personalID)
	if err != nil {
		return Employee{}, err
	}

	if err := es.validateEditData(email, phoneNumber, address); err != nil {
		return Employee{}, err
	}

	employee := employees[idx]
	employee.Email = email
	employee.PhoneNumber = phoneNumber
	employee.Address = address
	employees[idx] = employee

	if err := es.storage.Save(employees); err != nil {
		return Employee{}, err
	}

	return employee, nil
}
