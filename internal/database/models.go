// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package database

import (
	"time"
)

type Account struct {
	ID        int64
	FullName  string
	Balance   int64
	Currency  string
	CreatedAt time.Time
}

type Entry struct {
	ID        int64
	AccountID int64
	Amount    int64
	CreatedAt time.Time
}
