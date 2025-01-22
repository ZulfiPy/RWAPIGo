package vehicle

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/ZulfiPy/RWAPIGo/internal/storage"
)

type Vehicle struct {
	PlateNumber string
	Make        string
	Model       string
	Year        int
	FuelType    string
	Gearbox     string
	Color       string
	Body        string
}

type Vehicles map[string]Vehicle

type VehicleStorage struct {
	storage *storage.Storage[Vehicles]
}

func NewVehicleStorage(fileName string) *VehicleStorage {
	return &VehicleStorage{
		storage: storage.NewStorage[Vehicles](fileName),
	}
}

func (vs *VehicleStorage) GetStorage() *storage.Storage[Vehicles] {
	return vs.storage
}

func (vs *VehicleStorage) validateVehicle(input Vehicle) error {
	var fuelType = []string{"Petrol", "Diesel", "Hybrid", "Electric", "Lpg", "Cng"}
	var gearbox = []string{"Automatic", "Manual"}
	var colors = []string{"White", "Black", "Red", "Blue", "Green", "Yellow", "Gray", "Silver", "Brown"}
	var bodies = []string{"Sedan", "Touring", "Hatchback", "Minivan", "Coupe", "Cabriolet", "Pickup", "Limousine"}

	caser := cases.Title(language.English)

	if input.PlateNumber == "" {
		return errors.New("invalid input: vehicle plate number may not be empty")
	}

	if input.Make == "" {
		return errors.New("invalid input: vehicle make may not be empty")
	}

	if input.Model == "" {
		return errors.New("invalid input: vehicle model may not be empty")
	}

	if input.Year < 2010 || input.Year > time.Now().Year() {
		return errors.New("invalid input: vehicle year may not be lower than 2010 or greater than the current year")
	}

	if input.FuelType == "" {
		return errors.New("invalid input: vehicle fuel type may not be empty")
	}

	if !(slices.Contains(fuelType, caser.String(input.FuelType))) {
		return errors.New("invalid input: vehicle fuel type may only be (Petrol / Diesel / Hybrid / Electric / LPG / CNG)")
	}

	if input.Gearbox == "" {
		return errors.New("invalid input: vehicle gearbox may not be empty")
	}

	if !(slices.Contains(gearbox, caser.String(input.Gearbox))) {
		return errors.New("invalid input: vehicle gearbox may only be (Automatic or Manual)")
	}

	if input.Color == "" {
		return errors.New("invalid input: vehicle color may not be empty")
	}

	if !(slices.Contains(colors, caser.String(input.Color))) {
		return errors.New("invalid input: wrong vehicle color")
	}

	if input.Body == "" {
		return errors.New("invalid input: vehicle body may not be empty")
	}

	if !(slices.Contains(bodies, caser.String(input.Body))) {
		return errors.New("invalid input: wrong vehicle body")
	}

	return nil
}

func (vs *VehicleStorage) GetVehicle() (Vehicles, error) {
	vehicles := Vehicles{}

	if err := vs.storage.Load(&vehicles); err != nil {
		return nil, err
	}

	return vehicles, nil
}

func (vs *VehicleStorage) AddVehicle(input Vehicle) (Vehicle, error) {
	vehicles := Vehicles{}
	vs.storage.Load(&vehicles)

	if err := vs.storage.Load(&vehicles); err != nil {
		return Vehicle{}, err
	}

	if err := vs.validateVehicle(input); err != nil {
		return Vehicle{}, err
	}

	_, ok := vehicles[input.PlateNumber]

	if ok {
		return Vehicle{}, fmt.Errorf("vehiche with plate number %v is already in the storage", input.PlateNumber)
	}

	vehicles[input.PlateNumber] = input

	if err := vs.storage.Save(vehicles); err != nil {
		return Vehicle{}, err
	}

	return input, nil
}

func (vs *VehicleStorage) DeleteVehicle(plateNumber string) error {
	vehicles := Vehicles{}

	if err := vs.storage.Load(&vehicles); err != nil {
		return err
	}

	if _, ok := vehicles[plateNumber]; !ok {
		return fmt.Errorf("vehicle with plate number %v not found in the storage", plateNumber)
	}

	delete(vehicles, plateNumber)

	if err := vs.storage.Save(vehicles); err != nil {
		return err
	}

	return nil
}

func printFields(vehicle Vehicle) error {
	v := reflect.ValueOf(vehicle)
	typeOfVehicle := v.Type()

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Interface() == "" {
			return fmt.Errorf("%s cannot be empty", typeOfVehicle.Field(i).Name)
		}
	}
	return nil
}

func (vs *VehicleStorage) EditVehicle(input Vehicle) (Vehicle, error) {
	vehicles := Vehicles{}

	if err := vs.storage.Load(&vehicles); err != nil {
		return input, err
	}

	_, ok := vehicles[input.PlateNumber]

	if !ok {
		return input, fmt.Errorf("vehicle with plate number %v not found in the storage", input.PlateNumber)
	}

	if err := printFields(input); err != nil {
		return input, err
	}

	if vehicles[input.PlateNumber] == input {
		return Vehicle{}, errors.New("new data not detected")
	}

	vehicles[input.PlateNumber] = input

	if err := vs.storage.Save(vehicles); err != nil {
		return Vehicle{}, err
	}

	return input, nil
}
