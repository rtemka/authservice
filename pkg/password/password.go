// package password, used for hashing passwords
// and verifying provided plain passwords against hashed ones.
package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Hash generate hashed password using bcrypt, note
// that it'l work for password of any length including len == 0.
func Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to create hash password: %w", err)
	}
	return string(hash), nil
}

// Compare verifies that plain == hashed.
func Compare(plain, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}
