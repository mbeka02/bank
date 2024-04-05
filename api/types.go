package api

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

type CreateUserResponse struct {
	Username string `json:"user_name" validate:"required"`
	Fullname string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}
