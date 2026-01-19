package auth

import (
	"tfetch/model"
)

type AuthService struct{}

type AuthMethods interface {
	CallAuth() (*model.Token, error)
	CheckLogin(statusLogger model.StatusLoggerMethods) (*model.Token, error)
	RefreshToken(refreshToken string) (*model.Token, error)
	NewTokenStore(code string) (*TokenStore, error)
}

type Config struct {
	ClientId     string
	ClientSecret string
	Audience     string
	RedirectUri  string
	Scope        string
	Passphrase   string
}

type TokenStore struct {
	key  []byte
	salt []byte
}

type EncryptStore struct {
	Salt []byte `json:"salt"`
	IV   []byte `json:"iv"`
	Data []byte `json:"data"`
}

var (
	StateStore string
)
