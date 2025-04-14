package auth

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword_Success(t *testing.T) {
	password := "mysecurepassword"

	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if hashed == password {
		t.Errorf("hashed password should not match original password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil {
		t.Errorf("bcrypt failed to compare hash: %v", err)
	}
}

func TestCheckPasswordHash_Valid(t *testing.T) {
	password := "anothersecurepassword"
	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	err = CheckPasswordHash(hashed, password)
	if err != nil {
		t.Errorf("expected password to match, got error: %v", err)
	}
}

func TestCheckPasswordHash_Invalid(t *testing.T) {
	password := "correctpassword"
	hashed, _ := HashPassword(password)

	wrongPassword := "wrongpassword"
	err := CheckPasswordHash(hashed, wrongPassword)
	if err == nil {
		t.Errorf("expected error for wrong password, got nil")
	}
}
