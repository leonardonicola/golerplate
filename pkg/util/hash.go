package util

import "golang.org/x/crypto/bcrypt"

// The higher the salt, the more it takes to hash
// e.g: From 50ms (10) to 400ms (13)
const bcryptSalt = 10

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptSalt)
	return string(bytes), err
}

func CheckPasswordEquality(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
