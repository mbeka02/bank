package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mbeka02/bank/internal/auth"
	"github.com/mbeka02/bank/internal/database"
	"github.com/mbeka02/bank/utils"
	"log"
	"log/slog"
	"net/http"
	"strconv"
)

type APIServer struct {
	addr       string
	store      *database.Store
	tokenMaker auth.Maker
	config     utils.Config
}

type APIFunc func(w http.ResponseWriter, r *http.Request) error

type APIError struct {
	statusCode int
	message    string
}

func NewServer(addr string, store *database.Store, config utils.Config) (*APIServer, error) {
	maker, err := auth.NewJWTMaker(config.SymmetricKey)

	if err != nil {
		return nil, fmt.Errorf("Unable to create a new server instance: %v", err)
	}

	return &APIServer{
		addr:       addr,
		store:      store,
		tokenMaker: maker,
	}, nil
}

func (s *APIServer) Run() {
	router := http.NewServeMux()
	router.HandleFunc("GET /", modifyAPIFunc(s.handleGreetings))

	router.HandleFunc("POST /register", modifyAPIFunc(s.handleCreateUser))
	router.HandleFunc("POST /login", modifyAPIFunc(s.handleLogin))
	router.HandleFunc("GET /accounts", modifyAPIFunc(s.handleGetAccounts))
	router.HandleFunc("POST /accounts", modifyAPIFunc(s.handleCreateAccount))

	router.HandleFunc("GET /transfers", authMiddleware(modifyAPIFunc(s.handleGetTranfers), s.tokenMaker))
	router.HandleFunc("POST /transfers", modifyAPIFunc(s.handleTransferRequest))

	router.HandleFunc("GET /entries", modifyAPIFunc(s.handleGetEntries))

	router.HandleFunc("GET /accounts/{id}", modifyAPIFunc(s.handleGetAccount))
	router.HandleFunc("GET /entries/{id}", modifyAPIFunc(s.handleGetEntry))
	router.HandleFunc("GET /transfers/{id}", modifyAPIFunc(s.handleGetTransfer))
	log.Printf("Server is running on port %v", s.addr)

	err := http.ListenAndServe(s.addr, router)

	if err != nil {
		log.Fatal("Unable to spin up the server")
	}
}
func (e APIError) Error() string {
	return e.message
}

// handle JSON responses
func JSONResponse(w http.ResponseWriter, statusCode int, payload interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(&payload)
}

// returns a normal http handler function
func modifyAPIFunc(fn APIFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			if e, ok := err.(APIError); ok {
				JSONResponse(w, e.statusCode, e.message)

				slog.Error("API error", "err", e, "status", e.statusCode)
			}
		}
	}

}

// get path value and convert it to int64
func getIDFromRequest(r *http.Request) (int64, error) {
	id := r.PathValue("id")
	return strconv.ParseInt(id, 10, 64)
}

func (s *APIServer) validAccount(ctx context.Context, accountID int64, currency string) bool {
	acc, err := s.store.GetAccount(ctx, accountID)
	if err != nil {
		return false
	}
	if acc.Currency == currency {
		return true
	}
	return false

}
