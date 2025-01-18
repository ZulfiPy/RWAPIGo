package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ZulfiPy/RWAPIGo/internal/models/customer"
	"github.com/ZulfiPy/RWAPIGo/internal/models/vehicle"
	"github.com/ZulfiPy/RWAPIGo/internal/storage"
	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr    string
	customerStore *customer.CustomerStorage
	vehicleStore  *storage.Storage[vehicle.Vehicles]
}

func NewAPIServer(listenAddr string, customerStore *customer.CustomerStorage, vehicleStore *storage.Storage[vehicle.Vehicles]) *APIServer {
	return &APIServer{
		listenAddr:    listenAddr,
		customerStore: customerStore,
		vehicleStore:  vehicleStore,
	}
}

type APIError struct {
	Error string `json:"error"`
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	// router.HandleFunc("/", makeHTTPHandleFunc(s.homeHandler))

	// /customers
	router.HandleFunc("/customers", makeHTTPHandleFunc(s.handleCustomer))

	log.Println("JSON API server is running on port", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleCustomer(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.getCustomer(w, r)
	}
	if r.Method == "POST" {
		return s.addCustomer(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) getCustomer(w http.ResponseWriter, _ *http.Request) error {
	customers := customer.Customers{}
	err := s.customerStore.Load(&customers)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, customers)
}

func (s *APIServer) addCustomer(w http.ResponseWriter, r *http.Request) error {
	var newCustomer customer.Customer
	if err := json.NewDecoder(r.Body).Decode(&newCustomer); err != nil {
		return err
	}

	if err := s.customerStore.AddCustomer(newCustomer); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, newCustomer)
}

func WriteJSON(w http.ResponseWriter, status int, value any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(value)
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type apiError struct {
	Error string `json:"error"`
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, apiError{Error: err.Error()})
		}
	}
}
