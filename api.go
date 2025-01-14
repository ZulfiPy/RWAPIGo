package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr    string
	customerStore *Storage[Customers]
	vehicleStore  *Storage[Vehicles]
}

func NewAPIServer(listenAddr string, customerStore *Storage[Customers], vehicleStore *Storage[Vehicles]) *APIServer {
	return &APIServer{
		listenAddr:    listenAddr,
		customerStore: customerStore,
		vehicleStore:  vehicleStore,
	}
}

type Todo struct {
	userID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func (server *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/", makeHTTPHandleFunc(server.homeHandler))

	log.Println("JSON API server is running on port", server.listenAddr)

	http.ListenAndServe(server.listenAddr, router)
}

// home handler (sample)
func (server *APIServer) homeHandler(w http.ResponseWriter, r *http.Request) error {
	fmt.Println("Welcome to RWAPIGolang", w)
	fmt.Println("r", r)
	return nil
}

// test how API request works
// func (server *APIServer) getTodos(w http.ResponseWriter, r *http.Request) error {
// 	resp, err := http.Get("https://jsonplaceholder.typicode.com/todos")

// 	if err != nil {
// 		fmt.Println("Error occurred:", err)
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	var todos []Todo
// 	if err := json.NewDecoder(resp.Body).Decode(&todos); err != nil {
// 		http.Error(w, "Error decoding response", http.StatusInternalServerError)
// 		return err
// 	}

// 	return WriteJSON(w, http.StatusOK, todos)
// }

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
