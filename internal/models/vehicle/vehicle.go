package vehicle

import (
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
