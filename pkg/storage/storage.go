// package storage represents
// contract for working with database.
package storage

import "context"

// Storage is the DB contract
type Storage interface {
	AddUser(context.Context, User) error
	User(ctx context.Context, login string) (User, error)
}

// User
type User struct {
	Login string `json:"login"`
	Password
	CreatedAt  int64 `json:"created_at"`
	IsDisabled bool  `json:"is_disabled"`
}

// Password is the hashed password
// that stored in Storage
type Password struct {
	Hash        string `json:"hash"`
	GeneratedAt int64  `json:"generated_at"`
	IsActive    bool   `json:"is_active"`
}
