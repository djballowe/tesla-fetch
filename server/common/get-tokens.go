package common

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

var (
	stateStore string
	tokenStore = make(map[string]Token)
	storeMutex sync.Mutex
)

func GetTokens(writer http.ResponseWriter, req *http.Request) {
	var tokens Token
	err := json.NewDecoder(req.Body).Decode(&tokens)
	if err != nil {
		http.Error(writer, "Invalid request", http.StatusBadRequest)
		return
	}

	fmt.Println("Tokens: ", tokens)
	return
}

func GetTokenStore() (map[string]Token, string) {
	fmt.Println("Getting Token Store")
	storeMutex.Lock()
	defer storeMutex.Unlock()
	copyStore := make(map[string]Token)
	stateCopy := stateStore

	for k, v := range tokenStore {
		copyStore[k] = v
	}

	return copyStore, stateCopy
}
