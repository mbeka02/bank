package main

import (
	"fmt"
	"github.com/mbeka02/bank/internal/database"
	"log"
	"net/http"
)

type APIServer struct {
	Addr string
	DB   *database.Queries
}

func NewServer(addr string, queries *database.Queries) *APIServer {

	return &APIServer{
		Addr: addr,
		DB:   queries,
	}
}

func (s *APIServer) Run() {
	router := http.NewServeMux()

	router.HandleFunc("GET /", s.handleGreetings)
	router.HandleFunc("GET /accounts", s.handleGetAccounts)

	log.Printf("Server is running on port %v", s.Addr)

	err := http.ListenAndServe(s.Addr, router)

	if err != nil {
		log.Fatal("Unable to run the server")
	}

}

func (s *APIServer) handleGreetings(w http.ResponseWriter, r *http.Request) {
	name := "anthony"
	fmt.Fprintf(w, "Greetings : %v", name)

}

func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := s.DB.GetAccounts(r.Context())

	if err != nil {
		fmt.Println(err)
	}

	fmt.Fprintf(w, "%v", accounts)
}
