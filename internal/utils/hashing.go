package utils

import "golang.org/x/crypto/bcrypt"

const (
	// default cost for bcrypt hashing is 10, with 12 recommend for produciton use
	bCryptCost = 12
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bCryptCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CompareHashAndPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
