package auth

import(
	"net/http"
	"log"
	"errors"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	polkaKey := headers.Get("Authorization")
	if polkaKey == "" {
		err := errors.New("polkakey not found in headers")
		log.Println("\n",err)
		return "", err
	}

	words := strings.Split(polkaKey, " ")
	var clean []string
	for _, word := range words {
		if word != "ApiKey" && word != "Authorization:" {
			clean = append(clean, word)
		}
	}

	result := strings.Join(clean,"")

	return result, nil
}
