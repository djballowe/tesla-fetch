package auth

func GetTokenStore() (map[string]Token, string) {
	StoreMutex.Lock()
	defer StoreMutex.Unlock()
	copyStore := make(map[string]Token)
	stateCopy := StateStore

	// for k, v := TokenStore {
	// 	copyStore[k] = v
	// }

	return copyStore, stateCopy
}
