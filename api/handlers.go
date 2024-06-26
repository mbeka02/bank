package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
	"github.com/mbeka02/bank/internal/database"
	"github.com/mbeka02/bank/utils"
)

var validate *validator.Validate

func (s *APIServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	//get auth info from request context
	authPayload, err := getAuthPayload(r.Context())

	if err != nil {
		return err
	}
	accounts, err := s.store.GetAccounts(r.Context(), database.GetAccountsParams{
		Owner:  authPayload.Username,
		Limit:  30,
		Offset: 0,
	})
	if err != nil {
		return APIError{
			message:    "unable to process the request",
			statusCode: http.StatusInternalServerError,
		}
	}

	return JSONResponse(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetTranfers(w http.ResponseWriter, r *http.Request) error {
	transfers, err := s.store.GetTransfers(r.Context(), database.GetTransfersParams{
		Limit:  30,
		Offset: 0,
	})
	if err != nil {
		return APIError{
			message:    "unable to process the request",
			statusCode: http.StatusInternalServerError,
		}
	}
	return JSONResponse(w, http.StatusOK, transfers)
}

func (s *APIServer) handleGetEntries(w http.ResponseWriter, r *http.Request) error {
	entries, err := s.store.GetEntries(r.Context(), database.GetEntriesParams{
		Limit:  30,
		Offset: 0,
	})
	if err != nil {
		return APIError{
			message:    "unable to process the request",
			statusCode: http.StatusInternalServerError,
		}
	}

	return JSONResponse(w, http.StatusOK, entries)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	//get auth info from request context
	authPayload, err := getAuthPayload(r.Context())

	accountNum, err := getIDFromRequest(r)
	if err != nil {
		return err
	}
	account, err := s.store.GetAccount(r.Context(), accountNum)

	if err != nil {
		return APIError{
			message:    "unable to process the request",
			statusCode: http.StatusInternalServerError,
		}
	}
	if account.Owner != authPayload.Username {
		return APIError{
			message:    "account doesn't belong to the current user",
			statusCode: http.StatusUnauthorized,
		}
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
		return APIError{
			message:    "unable to process the request",
			statusCode: http.StatusInternalServerError,
		}
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
		return APIError{
			message:    "unable to process the request",
			statusCode: http.StatusInternalServerError,
		}

	}
	return JSONResponse(w, http.StatusOK, transfer)

}

func (s *APIServer) handleTransferRequest(w http.ResponseWriter, r *http.Request) error {
	//get auth info from request context
	authPayload, err := getAuthPayload(r.Context())
	if err != nil {
		return err
	}

	params := TransferTxRequest{}
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		return APIError{
			message:    "unable to process the request",
			statusCode: http.StatusInternalServerError,
		}
	}
	validate = validator.New()
	validate.RegisterValidation("currency", utils.ValidCurrency)
	if err := validate.Struct(params); err != nil {
		return APIError{
			message:    "field validation error" + err.Error(),
			statusCode: http.StatusBadRequest,
		}
	}

	senderAcc, valid := s.validAccount(r.Context(), params.SenderID, params.Currency)
	if !valid {
		return APIError{
			message:    "Invalid transfer details:transfer currency mismatch",
			statusCode: http.StatusBadRequest,
		}
	}
	if _, valid := s.validAccount(r.Context(), params.ReceiverID, params.Currency); !valid {
		return APIError{
			message:    "Invalid transfer details:transfer currency mismatch",
			statusCode: http.StatusBadRequest,
		}
	}
	if authPayload.Username != senderAcc.Owner {
		return APIError{
			message:    "sender account doesn't belong to the current user",
			statusCode: http.StatusUnauthorized,
		}
	}
	transferTxResult, err := s.store.TransferTx(r.Context(), database.TransferTxParams{
		SenderID:   params.SenderID,
		ReceiverID: params.ReceiverID,
		Amount:     params.Amount,
	})

	if err != nil {
		return APIError{
			message:    "unable to process the request",
			statusCode: http.StatusInternalServerError,
		}
	}
	return JSONResponse(w, http.StatusOK, transferTxResult)
}
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {

	//get auth info from request context
	authPayload, err := getAuthPayload(r.Context())

	if err != nil {
		return err
	}

	params := CreateAccountRequest{}
	err = json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		return APIError{
			message:    "unable to process the request",
			statusCode: http.StatusInternalServerError,
		}
	}
	validate = validator.New()
	if err := validate.Struct(params); err != nil {
		return APIError{
			message:    "field validation error:" + err.Error(),
			statusCode: http.StatusBadRequest,
		}
	}

	account, err := s.store.CreateAccount(r.Context(), database.CreateAccountParams{
		// a logged in use can only create an account for themselves
		Owner:    authPayload.Username,
		Currency: params.Currency,
		Balance:  0,
	})

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				return APIError{
					message:    "forbidden",
					statusCode: http.StatusForbidden,
				}
			}
		}
		return APIError{
			message:    "unable to process the request",
			statusCode: http.StatusInternalServerError,
		}
	}
	return JSONResponse(w, http.StatusOK, account)
}

func (s *APIServer) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	params := CreateUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		return APIError{
			message:    "unable to process the request",
			statusCode: http.StatusInternalServerError,
		}
	}
	validate = validator.New()
	if err := validate.Struct(params); err != nil {
		return APIError{
			message:    "field validation error" + err.Error(),
			statusCode: http.StatusBadRequest,
		}
	}
	passwordHash, err := utils.HashPassword(params.Password)
	if err != nil {
		return APIError{
			message:    err.Error(),
			statusCode: http.StatusInternalServerError,
		}
	}
	user, err := s.store.CreateUser(r.Context(), database.CreateUserParams{
		UserName: params.Username,
		FullName: params.Fullname,
		Email:    params.Email,
		Password: passwordHash,
	})

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				return APIError{
					message:    "forbidden: the username or email are already in use",
					statusCode: http.StatusForbidden,
				}
			}
		}
	}
	userResponse := newUserResponse(user)
	return JSONResponse(w, http.StatusOK, userResponse)

}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	params := LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&params)

	if err != nil {
		return APIError{
			message:    "unable to process the request",
			statusCode: http.StatusInternalServerError,
		}
	}
	validate = validator.New()
	if err := validate.Struct(params); err != nil {
		return APIError{
			message:    "field validation error" + err.Error(),
			statusCode: http.StatusBadRequest,
		}
	}
	user, err := s.store.GetUser(r.Context(), params.Username)
	if err != nil {
		return APIError{
			message:    "invalid username",
			statusCode: http.StatusNotFound,
		}
	}
	err = utils.ComparePassword(params.Password, user.Password)
	if err != nil {
		return APIError{
			message:    "unauthorized",
			statusCode: http.StatusUnauthorized,
		}
	}
	token, err := s.tokenMaker.CreateToken(user.UserName, s.config.AccessTokenDuration)
	userResponse := newUserResponse(user)
	rsp := LoginResponse{
		User:        userResponse,
		AccessToken: token,
	}
	return JSONResponse(w, http.StatusCreated, rsp)
}
