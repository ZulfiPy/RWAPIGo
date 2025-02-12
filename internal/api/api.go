package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ZulfiPy/RWAPIGo/internal/models/customer"
	"github.com/ZulfiPy/RWAPIGo/internal/models/employee"
	"github.com/ZulfiPy/RWAPIGo/internal/models/vehicle"
	"github.com/ZulfiPy/RWAPIGo/internal/utils"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr      string
	customerStorage *customer.CustomerStorage
	vehicleStorage  *vehicle.VehicleStorage
	employeeStorage *employee.EmployeeStorage
}

func NewAPIServer(listenAddr string, customerStorage *customer.CustomerStorage, vehicleStorage *vehicle.VehicleStorage, employeeStorage *employee.EmployeeStorage) *APIServer {
	return &APIServer{
		listenAddr:      listenAddr,
		customerStorage: customerStorage,
		vehicleStorage:  vehicleStorage,
		employeeStorage: employeeStorage,
	}
}

type APIError struct {
	Error string `json:"error"`
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	// /customers
	router.HandleFunc("/customers", makeHTTPHandleFunc(s.handleCustomer))
	router.HandleFunc("/customers/{personalID}/vehicles", makeHTTPHandleFunc(s.handleCustomerVehicle))
	router.HandleFunc("/customers/{personalID}/{plateNumber}/delete-vehicle", makeHTTPHandleFunc(s.handleDeleteVehicleFromCustomer))

	router.HandleFunc("/vehicles", makeHTTPHandleFunc(s.handleVehicle))

	router.HandleFunc("/employees", makeHTTPHandleFunc(s.handleEmployee))

	log.Println("JSON API server is running on port", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleCustomer(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetCustomer(w, r)
	}
	if r.Method == "POST" {
		return s.handleAddCustomer(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteCustomer(w, r)
	}
	if r.Method == "PUT" {
		return s.handleEditCustomer(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleVehicle(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetVehicle(w, r)
	}
	if r.Method == "POST" {
		return s.handleAddVehicle(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteVehicle(w, r)
	}
	if r.Method == "PUT" {
		return s.handleEditVehicle(w, r)
	}
	return fmt.Errorf("method %s not allowed", r.Method)
}

func (s *APIServer) handleEmployee(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetEmployee(w, r)
	}
	if r.Method == "POST" {
		return s.handleAddEmployee(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteEmployee(w, r)
	}
	if r.Method == "PUT" {
		return s.handleEditEmployee(w, r)
	}
	return fmt.Errorf("method %s not allowed", r.Method)
}

func (s *APIServer) handleCustomerVehicle(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		return s.handleAddVehicleToCustomer(w, r)
	}
	return fmt.Errorf("method %s not allowed", r.Method)
}

func (s *APIServer) handleGetCustomer(w http.ResponseWriter, _ *http.Request) error {
	customers, err := s.customerStorage.GetCustomers()

	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
	}

	return WriteJSON(w, http.StatusOK, customers)

}

func (s *APIServer) handleGetVehicle(w http.ResponseWriter, _ *http.Request) error {
	vehicles, err := s.vehicleStorage.GetVehicles()

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusAccepted, vehicles)
}

func (s *APIServer) handleGetEmployee(w http.ResponseWriter, _ *http.Request) error {
	employees, err := s.employeeStorage.GetEmployees()

	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
	}

	return WriteJSON(w, http.StatusOK, employees)
}

func (s *APIServer) handleAddCustomer(w http.ResponseWriter, r *http.Request) error {
	var newCustomer customer.Customer
	if err := json.NewDecoder(r.Body).Decode(&newCustomer); err != nil {
		return err
	}

	if err := s.customerStorage.AddCustomer(newCustomer); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, newCustomer)
}

func (s *APIServer) handleAddVehicle(w http.ResponseWriter, r *http.Request) error {
	var newVehicle vehicle.Vehicle
	if err := json.NewDecoder(r.Body).Decode(&newVehicle); err != nil {
		return err
	}

	vehicle, err := s.vehicleStorage.AddVehicle(newVehicle)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
	}

	return WriteJSON(w, http.StatusOK, vehicle)
}

func (s *APIServer) handleAddEmployee(w http.ResponseWriter, r *http.Request) error {
	var newEmployee employee.Employee

	if err := json.NewDecoder(r.Body).Decode(&newEmployee); err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
	}

	employee, err := s.employeeStorage.AddEmployee(newEmployee)

	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
	}

	return WriteJSON(w, http.StatusOK, employee)
}

type CustomResponse struct {
	Response string `json:"response"`
}

func (s *APIServer) handleDeleteCustomer(w http.ResponseWriter, r *http.Request) error {
	var personalID struct {
		PersonalID int64 `json:"PersonalID"`
	}
	if err := json.NewDecoder(r.Body).Decode(&personalID); err != nil {
		return err
	}

	personalIDLength := utils.IntLength(personalID.PersonalID)

	if personalIDLength != 11 {
		return WriteJSON(w, http.StatusBadRequest, ApiError{Error: "personalID must be exactly 11 digits"})
	}

	if err := s.customerStorage.DeleteCustomer(personalID.PersonalID); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, CustomResponse{Response: "customer deleted"})
}

func (s *APIServer) handleDeleteVehicle(w http.ResponseWriter, r *http.Request) error {
	var plateNumber struct {
		PlateNumber string `json:"PlateNumber"`
	}

	if err := json.NewDecoder(r.Body).Decode(&plateNumber); err != nil {
		return err
	}

	if err := s.vehicleStorage.DeleteVehicle(plateNumber.PlateNumber); err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
	}

	return WriteJSON(w, http.StatusOK, CustomResponse{Response: "vehicle deleted"})
}

func (s *APIServer) handleDeleteVehicleFromCustomer(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	personalID, err := strconv.ParseInt(vars["personalID"], 10, 64)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
	}

	plateNumber := vars["plateNumber"]
	if err := s.vehicleStorage.GetVehicle(plateNumber); err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
	}

	if err := s.customerStorage.DeleteVehicle(plateNumber, personalID); err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
	}

	return WriteJSON(w, http.StatusOK, CustomResponse{Response: "vehicle deleted from customer"})
}

func (s *APIServer) handleDeleteEmployee(w http.ResponseWriter, r *http.Request) error {
	var personalID struct {
		PersonalID int64 `json:"PersonalID"`
	}
	if err := json.NewDecoder(r.Body).Decode(&personalID); err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
	}

	if err := s.employeeStorage.DeleteEmployee(personalID.PersonalID); err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
	}
	return nil
}

func (s *APIServer) handleEditCustomer(w http.ResponseWriter, r *http.Request) error {
	var editData struct {
		FirstName   string `json:"FirstName"`
		LastName    string `json:"LastName"`
		Email       string `json:"Email"`
		PhoneNumber string `json:"PhoneNumber"`
		PersonalID  int64  `json:"PersonalID"`
	}

	if err := json.NewDecoder(r.Body).Decode(&editData); err != nil {
		return err
	}

	if err := s.customerStorage.EditCustomer(editData.LastName, editData.FirstName, editData.Email, editData.PhoneNumber, editData.PersonalID); err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
	}

	return WriteJSON(w, http.StatusOK, CustomResponse{Response: "customer successfully edited"})
}

func (s *APIServer) handleEditVehicle(w http.ResponseWriter, r *http.Request) error {
	var editVehicle vehicle.Vehicle
	if err := json.NewDecoder(r.Body).Decode(&editVehicle); err != nil {
		return err
	}
	vehicle, err := s.vehicleStorage.EditVehicle(editVehicle)

	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
	}

	return WriteJSON(w, http.StatusOK, vehicle)
}

func (s *APIServer) handleEditEmployee(w http.ResponseWriter, r *http.Request) error {
	var editCustomerData struct {
		PersonalID  int64  `json:"PersonalID"`
		Email       string `json:"Email"`
		PhoneNumber string `json:"PhoneNumber"`
		Address     string `json:"Address"`
	}

	if err := json.NewDecoder(r.Body).Decode(&editCustomerData); err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
	}

	employee, err := s.employeeStorage.EditEmployeeContacts(editCustomerData.Email, editCustomerData.PhoneNumber, editCustomerData.Address, editCustomerData.PersonalID)

	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
	}

	return WriteJSON(w, http.StatusOK, employee)
}

func (s *APIServer) handleAddVehicleToCustomer(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	personalID, personalIDErr := strconv.ParseInt(vars["personalID"], 10, 64)
	if personalIDErr != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: personalIDErr.Error()})
	}

	var vehicle vehicle.Vehicle
	if err := json.NewDecoder(r.Body).Decode(&vehicle); err != nil {
		return err
	}

	if err := s.vehicleStorage.GetVehicle(vehicle.PlateNumber); err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
	}

	customer, err := s.customerStorage.AddVehicle(vehicle, personalID)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
	}

	return WriteJSON(w, http.StatusOK, customer)
}

type ApiFunc func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, value any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(value)
}

func makeHTTPHandleFunc(f ApiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}
