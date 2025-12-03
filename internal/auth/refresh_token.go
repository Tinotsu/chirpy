package auth

import (
	"crypto/rand"
	"log"
	"encoding/hex"
)

func MakeRefreshToken() (string, error) {
	buf := make([]byte, 32)
	n, err := rand.Read(buf)
	if err != nil {
		log.Println("\n",err)
		return "", err
	}
	buf = append(buf, byte(n))
	encoded := hex.EncodeToString(buf)
	return encoded, nil
}
