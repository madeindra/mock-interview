package util

import (
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

func GenerateRandom() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 10

	random := make([]byte, length)
	for i := range random {
		random[i] = charset[rand.Intn(len(charset))]
	}

	return string(random)
}

func CreateHash(plain string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CompareHash(plain, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}
