// package domain represents core types
// for the authentication service.
package domain

import "context"

// User - user's entity.
type User struct {
	Login string `json:"login"`
	Password
	CreatedAt  int64 `json:"created_at"`
	IsDisabled bool  `json:"is_disabled"`
}

// Password is the hashed password.
type Password struct {
	Hash        string `json:"-"`
	GeneratedAt int64  `json:"generated_at"`
	IsActive    bool   `json:"is_active"`
}

// Storage is the contract that database must implement.
type Storage interface {
	AddUser(context.Context, User) error
	User(ctx context.Context, login string) (User, error)
	UpdatePassword(context.Context, User) error
	Password(context.Context, User) (Password, error)
}
