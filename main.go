package main

import (
	"fmt"
	"os"
)

func ensureStorageFile[T any](storage *Storage[T], data T) error {
	dataPath := fmt.Sprintf("./%s", storage.FileName)
	
	if _, err := os.Stat(dataPath); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("creating file", storage.FileName)
			storage.Save(data)
			return nil
		}
		fmt.Println("Error accessing file:", err)
		return err
	}

	fmt.Printf("file exists %s, no action needed.\n", storage.FileName)
	return nil
}


func main() {
	fmt.Println("RWAPIGolang runs...")
	customers := Customers{}
	vehicles := Vehicles{}

	customerStorage := NewStorage[Customers]("customers.json")
	ensureStorageFile(customerStorage, customers)

	vehicleStorage := NewStorage[Vehicles]("vehicles.json")
	ensureStorageFile(vehicleStorage, vehicles)
}
