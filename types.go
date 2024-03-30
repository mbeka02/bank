package main

type CreateAccountRequest struct {
	FullName string `json:"full_name" validate:"required"`
	Currency string `json:"currency" validate:"required"`
}

type TransferTxRequest struct {
	Amount     int64 `json:"amount" validate:"required"`
	SenderID   int64 `json:"sender_id" validate:"required"`
	ReceiverID int64 `json:"receiver_id" validate:"required"`
}
