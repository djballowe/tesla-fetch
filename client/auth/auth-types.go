package auth

import "sync"

type AuthResponse struct {
	CallbackUrl string `json:"callback_url"`
}

type Config struct {
	ClientId     string
	ClientSecret string
	Audience     string
	RedirectUri  string
	Scope        string
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
	State        string `json:"state"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type TokenStore struct {
	filePath string
	key      []byte
}

type EncryptStore struct {
	Salt []byte `json:"salt"`
	IV   []byte `json:"iv"`
	Data []byte `json:"data"`
}

var (
	StateStore string
	StoreMutex sync.Mutex
)
