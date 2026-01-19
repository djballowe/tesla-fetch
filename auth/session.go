package auth

import (
	"os"
	"tfetch/model"
)

func (a *AuthService) CheckLogin(status model.StatusLoggerMethods) (*model.Token, error) {
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
