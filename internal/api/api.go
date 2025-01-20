package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ZulfiPy/RWAPIGo/internal/models/customer"
	"github.com/ZulfiPy/RWAPIGo/internal/models/vehicle"
	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr      string
	customerStorage *customer.CustomerStorage
	vehicleStorage  *vehicle.VehicleStorage
}

func NewAPIServer(listenAddr string, customerStorage *customer.CustomerStorage, vehicleStorage *vehicle.VehicleStorage) *APIServer {
	return &APIServer{
		listenAddr:      listenAddr,
		customerStorage: customerStorage,
		vehicleStorage:  vehicleStorage,
	}
}

type APIError struct {
	Error string `json:"error"`
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	// /customers
	router.HandleFunc("/customers", makeHTTPHandleFunc(s.handleCustomer))

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

func (s *APIServer) handleGetCustomer(w http.ResponseWriter, _ *http.Request) error {
	customers, err := s.customerStorage.GetCustomer()

	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
	}

	return WriteJSON(w, http.StatusOK, customers)

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

	personalIDLength := customer.IntLength(personalID.PersonalID)

	if personalIDLength != 11 {
		return WriteJSON(w, http.StatusBadRequest, ApiError{Error: "personalID must be exactly 11 digits"})
	}

	if err := s.customerStorage.DeleteCustomer(personalID.PersonalID); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, CustomResponse{Response: "customer deleted"})
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
