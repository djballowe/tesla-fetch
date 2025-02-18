package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
)

func NewTokeStore(code string) (*TokenStore, error) {
	if code == "" {
		return nil, errors.New("missing or invalid passcode")
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}

	tfetchDir := filepath.Join(configDir, "tfetch")
	err = os.MkdirAll(tfetchDir, 0700)
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(tfetchDir, "tfetch-tokens.dat")

	key, err := CreateKey([]byte(code), nil)
	if err != nil {
		return nil, err
	}

	return &TokenStore{
		filePath: filePath,
		key:      key,
	}, nil
}

func CreateKey(code []byte, salt []byte) ([]byte, error) {
	if salt == nil {
		salt = make([]byte, 16)
		_, err := io.ReadFull(rand.Reader, salt)
		if err != nil {
			return nil, errors.New("Could not create token key")
		}
	}

	k := sha256.New()
	k.Write(salt)
	k.Write(code)
	return k.Sum(nil), nil
}

func (store *TokenStore) SaveTokens(tokens *Token) error {
	data, err := json.Marshal(tokens)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(store.key)
	if err != nil {
		return err
	}

	iv := make([]byte, 12)
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	encrypt := gcm.Seal(nil, iv, data, nil)

	encryptData, err := json.Marshal(encrypt)
	if err != nil {
		return err
	}

	return os.WriteFile(store.filePath, encryptData, 0600)
}
