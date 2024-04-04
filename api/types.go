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
