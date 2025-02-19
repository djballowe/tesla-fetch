package auth

import (
	"sync"
	"time"
)

type AuthResponse struct {
	CallbackUrl string `json:"callback_url"`
}

type Config struct {
	ClientId     string
	ClientSecret string
	Audience     string
	RedirectUri  string
	Scope        string
	Passphrase   string
}

type Token struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	IdToken      string    `json:"id_token"`
	State        string    `json:"state"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	CreateAt     time.Time `json:"created_at"`
}

type TokenStore struct {
	filePath string
	key      []byte
	salt     []byte
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
