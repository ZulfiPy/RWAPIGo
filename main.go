package main

import (
	"fmt"
)

func main() {
	fmt.Println("RWAPIGolang runs...")
	customers := Customers{}
	vehicles := Vehicles{}

	customerStorage := NewStorage[Customers]("customers.json")
	EnsureStorageFile(customerStorage, customers)

	vehicleStorage := NewStorage[Vehicles]("vehicles.json")
	EnsureStorageFile(vehicleStorage, vehicles)

	server := NewAPIServer(":8080", customerStorage, vehicleStorage)
	server.Run()
}
