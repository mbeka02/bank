package api

import "github.com/mbeka02/bank/internal/database"

type CreateAccountRequest struct {
	Owner    string `json:"owner" validate:"required"`
	Currency string `json:"currency" validate:"required"`
}

type TransferTxRequest struct {
	Amount     int64  `json:"amount" validate:"required,min=1"`
	SenderID   int64  `json:"sender_id" validate:"required,min=1"`
	ReceiverID int64  `json:"receiver_id" validate:"required,gt=0"`
	Currency   string `json:"currency" validate:"required,currency"`
}

type CreateUserRequest struct {
	Username string `json:"user_name" validate:"required"`
	Fullname string `json:"full_name" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
	Email    string `json:"email" validate:"required,email"`
}

type UserResponse struct {
	Username string `json:"user_name" validate:"required"`
	Fullname string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}

type LoginRequest struct {
	Username string `json:"user_name" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}

func newUserResponse(user database.User) UserResponse {
	return UserResponse{
		Username: user.UserName,
		Fullname: user.FullName,
		Email:    user.Email,
	}
}
