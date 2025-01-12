package main

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
