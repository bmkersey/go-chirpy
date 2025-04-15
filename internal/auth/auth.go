package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
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

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		Subject:   userID.String(),
	})

	signedToken, err := newToken.SignedString([]byte(tokenSecret))
	if err != nil {
		log.Printf("Error signing token: %s", err)
		return "", err
	}

	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(parsedToken *jwt.Token) (interface{}, error) {
		if _, ok := parsedToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", parsedToken.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	if !parsedToken.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if len(authHeader) < 1 {
		return "", fmt.Errorf("no authorization header found: %v", authHeader)
	}
	authHeader = strings.TrimSpace(strings.Replace(authHeader, "Bearer", "", -1))
	return authHeader, nil
}

func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)
	rand.Read(token)

	stringToken := hex.EncodeToString(token)

	return stringToken, nil
}
