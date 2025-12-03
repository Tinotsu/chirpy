// Package auth
package auth

import (
	"log"
	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	hashedPw, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Println("\n",err)
		return "", err
	}
	return hashedPw, nil
}

func CheckPasswordHash(password string, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}
