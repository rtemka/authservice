// package hash, used for hashing passwords
// and verifying provided plain passwords against hashed ones.
package hash

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword generate hashed password using bcrypt, note
// that it'l work for password of any length including len == 0.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to create hash password: %w", err)
	}
	return string(hash), nil
}

// ComparePasswords verifies that plain == hashed.
func ComparePasswords(plain, hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
}
