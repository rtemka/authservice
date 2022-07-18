// package storage represents
// contract for working with database.
package storage

import (
	"authservice/domain"
	"context"
)

// Storage is the DB contract
type Storage interface {
	AddUser(context.Context, domain.User) error
	User(ctx context.Context, login string) (domain.User, error)
	UpdatePassword(context.Context, domain.User) error
	Password(context.Context, domain.User) (domain.Password, error)
}
