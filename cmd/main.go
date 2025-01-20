package main

import (
	"fmt"

	"github.com/ZulfiPy/RWAPIGo/internal/api"
	"github.com/ZulfiPy/RWAPIGo/internal/models/customer"
	"github.com/ZulfiPy/RWAPIGo/internal/models/vehicle"
	"github.com/ZulfiPy/RWAPIGo/internal/storage"
)

func main() {
	fmt.Println("RWAPIGolang runs...")

	customers := customer.Customers{}
	vehicles := vehicle.Vehicles{}

	customerStorage := customer.NewCustomerStorage("customers.json")
	storage.EnsureStorageFile(customerStorage.GetStorage(), customers)

	vehicleStorage := vehicle.NewVehicleStorage("vehicles.json")
	storage.EnsureStorageFile(vehicleStorage.GetStorage(), vehicles)

	server := api.NewAPIServer(":8080", customerStorage, vehicleStorage)
	server.Run()
}
