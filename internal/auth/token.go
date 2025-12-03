// Package auth
package auth

import (
	"log"
	"strings"
	"errors"
	"github.com/google/uuid"
	"time"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

type TokenType string

const (
	TokenTypeAccess TokenType = "chirpy-access"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	signingKey := []byte(tokenSecret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    string(TokenTypeAccess),
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})
	return token.SignedString(signingKey)
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		log.Println("Validation, parsing: ",err)
		log.Println("token: ", token)
		log.Println("tokenstring: -", tokenString,"-")
		return uuid.Nil, err
	}

	userID, err := token.Claims.GetSubject()
	if err != nil {
		log.Println("\nValidation, userID: ",err)
		return uuid.Nil, err
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil {
		log.Println("\nValidation, issuer: ",err)
		return uuid.Nil, err
	}
	if issuer != string(TokenTypeAccess) {
		log.Println("\nValidation, issuer invalid: ",err)
		return uuid.Nil, errors.New("invalid issuer")
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		log.Println("\nValidation uuid parse",err)
		return uuid.Nil, err
	}

	return userUUID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	token := headers.Get("Authorization")
	if token == "" {
		err := errors.New("bearer token not found in headers")
		log.Println("\n",err)
		return "", err
	}

	words := strings.Split(token, " ")
	var clean []string
	for _, word := range words {
		if word != "Bearer" {
			clean = append(clean, word)
		}
	}

	result := strings.Join(clean,"")

	return result, nil
}
