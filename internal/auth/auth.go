package auth

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)


func HashPassword(password string)(string, error){
	hashed_pw, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		log.Printf("error hashing password: %s", err)
		return "", err
	}
	return string(hashed_pw), nil
}

func CheckPasswordHash(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}