package employee

import (
	"time"

	"github.com/ZulfiPy/RWAPIGo/internal/storage"
)

type Employee struct {
	FirstName   string
	LastName    string
	PersonalID  int64
	DateOfBirth time.Time
	Email       string
	PhoneNumber string
	Address     string
}

type Employees []Employee

type EmployeeStorage struct {
	storage *storage.Storage[Employees]
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