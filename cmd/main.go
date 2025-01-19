package main

import (
	"fmt"

	"github.com/ZulfiPy/RWAPIGo/internal/api"
	"github.com/ZulfiPy/RWAPIGo/internal/models/customer"
	"github.com/ZulfiPy/RWAPIGo/internal/storage"
)

func main() {
	fmt.Println("RWAPIGolang runs...")
	customers := customer.Customers{}

	customerStorage := customer.NewCustomerStorage("customers.json")
	storage.EnsureStorageFile(customerStorage.GetStorage(), customers)

	server := api.NewAPIServer(":8080", customerStorage)
	server.Run()
}
