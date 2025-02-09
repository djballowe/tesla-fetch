package helpers

import (
	"log"
	"tesla-app/client/auth"
)

func GetTokenStore() (map[string]auth.Token, string) {
	log.Println("Getting Token Store")
	auth.StoreMutex.Lock()
	defer auth.StoreMutex.Unlock()
	copyStore := make(map[string]auth.Token)
	stateCopy := auth.StateStore

	for k, v := range auth.TokenStore {
		copyStore[k] = v
	}

	return copyStore, stateCopy
}
