// package auth

// import (
// 	"fmt"

// 	"golang.org/x/crypto/bcrypt"
// )

// func HashPassword(password string) (string, error) {
// 	hashedPw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return "", fmt.Errorf("Error hashing password: %w", err)
// 	}

// 	return string(hashedPw), nil
// }

// /**
// 	Returns nil on success
// 	Returns error on error
// */
// func CheckPasswordHash(password, hash string) error {
// 	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
// 	if err != nil {
// 		return fmt.Errorf("passwords do not match: %w", err)
// 	}
// 	return nil
// }

package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the given password using bcrypt.
func HashPassword(password string) (string, error) {
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %w", err)
	}
	return string(hashedPw), nil
}

// CheckPasswordHash compares a plain password with its hashed version.
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}