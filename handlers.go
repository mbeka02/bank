package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mbeka02/bank/internal/database"
)

type APIServer struct {
	Addr string
	DB   *database.Queries
}

type APIFunc func(w http.ResponseWriter, r *http.Request) error

type APIError struct {
	Error string `json:"error"`
}

// handle JSON responses
func JSONResponse(w http.ResponseWriter, statusCode int, payload interface{}) error {

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(&payload)
}

// modify the function signature to that of a normal handler function
func modifyAPIFunc(fn APIFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			JSONResponse(w, http.StatusBadRequest, APIError{err.Error()})
		}
	}

}
func NewServer(addr string, queries *database.Queries) *APIServer {

	return &APIServer{
		Addr: addr,
		DB:   queries,
	}
}

func (s *APIServer) Run() {
	router := http.NewServeMux()

	router.HandleFunc("GET /", modifyAPIFunc(s.handleGreetings))
	router.HandleFunc("GET /accounts", modifyAPIFunc(s.handleGetAccounts))
	router.HandleFunc("GET /transactions", modifyAPIFunc(s.handleGetTransactions))

	log.Printf("Server is running on port %v", s.Addr)

	err := http.ListenAndServe(s.Addr, router)

	if err != nil {
		log.Fatal("Unable to spin up the server")
	}

}

func (s *APIServer) handleGreetings(w http.ResponseWriter, r *http.Request) error {
	name := "anthony"
	return JSONResponse(w, http.StatusOK, name)

}

func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.DB.GetAccounts(r.Context(), database.GetAccountsParams{
		Limit:  40,
		Offset: 2,
	})
	if err != nil {
		return err
	}

	return JSONResponse(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetTransactions(w http.ResponseWriter, r *http.Request) error {
	transactions, err := s.DB.GetTransactions(r.Context(), database.GetTransactionsParams{
		Limit:  5,
		Offset: 5,
	})
	if err != nil {
		return err
	}
	return JSONResponse(w, http.StatusOK, transactions)
}
