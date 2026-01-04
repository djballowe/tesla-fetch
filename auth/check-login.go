package auth

import (
	"os"
	"tfetch/ui"
)

func (a *AuthService) CheckLogin(status ui.StatusLoggerMethods) (*Token, error) {
	passphrase := os.Getenv("PASSPHRASE")
	store := TokenStore{}
	token, err := store.LoadTokens(passphrase)
	if err != nil {
		if err.Error() == "No token stored" {
			token, err = a.CallAuth()
			if err != nil {
				return nil, err
			}
		}
	}

	if store.IsExpired(token.CreateAt, token.ExpiresIn) {
		status.Log("Token expired, refreshing")
		token, err = a.RefreshToken(token.RefreshToken)
		if err != nil {
			return nil, err
		}
	}

	return token, nil
}
