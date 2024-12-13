package util

import "golang.org/x/crypto/bcrypt"

const bcryptSalt = 13

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptSalt)
	return string(bytes), err
}

func CheckPasswordEquality(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
