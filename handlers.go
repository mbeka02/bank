package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/mbeka02/bank/internal/database"
)

type APIServer struct {
	Addr  string
	store *database.Store
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

// returns a normal http handler function
func modifyAPIFunc(fn APIFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			JSONResponse(w, http.StatusBadRequest, APIError{err.Error()})
		}
	}

}
func NewServer(addr string, store *database.Store) *APIServer {

	return &APIServer{
		Addr:  addr,
		store: store,
	}
}

func (s *APIServer) Run() {
	router := http.NewServeMux()

	router.HandleFunc("GET /", modifyAPIFunc(s.handleGreetings))
	router.HandleFunc("GET /accounts", modifyAPIFunc(s.handleGetAccounts))
	router.HandleFunc("POST /accounts", modifyAPIFunc(s.handleCreateAccount))

	router.HandleFunc("GET /transfers", modifyAPIFunc(s.handleGetTranfers))
	router.HandleFunc("POST /transfers", modifyAPIFunc(s.handleTransferTx))

	router.HandleFunc("GET /entries", modifyAPIFunc(s.handleGetEntries))

	router.HandleFunc("GET /accounts/{id}", modifyAPIFunc(s.handleGetAccount))
	router.HandleFunc("GET /entries/{id}", modifyAPIFunc(s.handleGetEntry))
	router.HandleFunc("GET /transfers/{id}", modifyAPIFunc(s.handleGetTransfer))
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
	accounts, err := s.store.GetAccounts(r.Context(), database.GetAccountsParams{
		Limit:  30,
		Offset: 0,
	})
	if err != nil {
		return err
	}

	return JSONResponse(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetTranfers(w http.ResponseWriter, r *http.Request) error {
	transfers, err := s.store.GetTransfers(r.Context(), database.GetTransfersParams{
		Limit:  30,
		Offset: 0,
	})
	if err != nil {
		return err
	}
	return JSONResponse(w, http.StatusOK, transfers)
}

func (s *APIServer) handleGetEntries(w http.ResponseWriter, r *http.Request) error {
	entries, err := s.store.GetEntries(r.Context(), database.GetEntriesParams{
		Limit:  30,
		Offset: 0,
	})
	if err != nil {
		return err
	}

	return JSONResponse(w, http.StatusOK, entries)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	i, err := getIDFromRequest(r)
	if err != nil {
		return err
	}
	account, err := s.store.GetAccount(r.Context(), i)

	if err != nil {
		return err
	}
	return JSONResponse(w, http.StatusOK, account)
}

func (s *APIServer) handleGetEntry(w http.ResponseWriter, r *http.Request) error {
	i, err := getIDFromRequest(r)
	if err != nil {
		return err
	}
	entry, err := s.store.GetEntry(r.Context(), i)

	if err != nil {
		return err
	}
	return JSONResponse(w, http.StatusOK, entry)
}

func (s *APIServer) handleGetTransfer(w http.ResponseWriter, r *http.Request) error {
	i, err := getIDFromRequest(r)
	if err != nil {
		return err
	}
	transfer, err := s.store.GetEntry(r.Context(), i)

	if err != nil {
		return err
	}
	return JSONResponse(w, http.StatusOK, transfer)

}

func (s *APIServer) handleTransferTx(w http.ResponseWriter, r *http.Request) error {

	params := TransferTxRequest{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		return err
	}
	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		return err
	}
	transferTxResult, err := s.store.TransferTx(r.Context(), database.TransferTxParams{
		SenderID:   params.SenderID,
		ReceiverID: params.ReceiverID,
		Amount:     params.Amount,
	})

	if err != nil {
		return err
	}
	return JSONResponse(w, http.StatusOK, transferTxResult)
}
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	params := CreateAccountRequest{}
	err := json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		return err
	}
	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		return err
	}

	account, err := s.store.CreateAccount(r.Context(), database.CreateAccountParams{
		FullName: params.FullName,
		Currency: params.Currency,
		Balance:  0,
	})

	if err != nil {

		return err
	}
	return JSONResponse(w, http.StatusOK, account)
}

// get path value and convert it to int64
func getIDFromRequest(r *http.Request) (int64, error) {
	id := r.PathValue("id")
	return strconv.ParseInt(id, 10, 64)
}
