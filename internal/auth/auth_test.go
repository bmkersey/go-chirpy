package auth

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
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

func TestJWTCreationAndValidation(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret-key"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}
	if token == "" {
		t.Fatal("Expected non-empty token")
	}

	extractedID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}
	if extractedID != userID {
		t.Fatalf("Expected user ID %v, got %v", userID, extractedID)
	}
}

func TestJWTExpiration(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret-key"
	expiresIn := -time.Hour 

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	_, err = ValidateJWT(token, tokenSecret)
	if err == nil {
		t.Fatal("Expected error for expired token, got none")
	}
}

func TestJWTInvalidSecret(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret-key"
	wrongSecret := "wrong-secret-key"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Fatal("Expected error for token with wrong secret, got none")
	}
}

func TestJWTInvalidFormat(t *testing.T) {
	_, err := ValidateJWT("not-a-valid-token", "any-secret")
	if err == nil {
		t.Fatal("Expected error for malformed token, got none")
	}
}

func TestJWTManuallyTampered(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "test-secret-key"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		t.Fatal("Expected JWT to have 3 parts")
	}
	payload := parts[1]
	if len(payload) > 0 {
		chars := []rune(payload)
		if chars[0] == 'a' {
			chars[0] = 'b'
		} else {
			chars[0] = 'a'
		}
		parts[1] = string(chars)
	}
	
	tamperedToken := strings.Join(parts, ".")
	
	_, err = ValidateJWT(tamperedToken, tokenSecret)
	if err == nil {
		t.Fatal("Expected error for tampered token, got none")
	}
}

func TestGetBearerToken_Valid(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer abc.def.ghi")

	token, err := GetBearerToken(headers)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if token != "abc.def.ghi" {
		t.Errorf("expected token to be 'abc.def.ghi', got '%s'", token)
	}
}

func TestGetBearerToken_MissingHeader(t *testing.T) {
	headers := http.Header{}

	token, err := GetBearerToken(headers)
	if err == nil {
		t.Error("expected error for missing header, got nil")
	}
	if token != "" {
		t.Errorf("expected empty token, got '%s'", token)
	}
}

func TestGetBearerToken_EmptyBearer(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer ")

	token, err := GetBearerToken(headers)
	if err != nil {
		t.Errorf("expected no error for empty bearer, got %v", err)
	}
	if token != "" {
		t.Errorf("expected empty token string, got '%s'", token)
	}
}

func TestGetBearerToken_ExtraWhitespace(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer    abc.def.ghi    ")

	token, err := GetBearerToken(headers)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if token != "abc.def.ghi" {
		t.Errorf("expected token to be 'abc.def.ghi', got '%s'", token)
	}
}