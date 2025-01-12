package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("RWAPIGolang runs...")
	customers := Customers{}

	customersStorage := NewStorage[Customers]("customers.json")

	customersPath := fmt.Sprintf("./%s", customersStorage.FileName)

	_, err := os.Stat(customersPath)

	if os.IsNotExist(err) {
		fmt.Println("creating", customersStorage.FileName)
		customersStorage.Save(customers)
	} else if err != nil {
		fmt.Println("Error accessing file:", err)
	} else {
		fmt.Println(customersStorage.FileName, "file exists, no action needed.")
	}
}
