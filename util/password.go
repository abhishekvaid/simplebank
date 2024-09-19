package util

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func CheckPassword(plain string, hashed string) (err error) {
	err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	return
}
