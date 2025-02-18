package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"io"
	"log"
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

	key, salt, err := CreateKey([]byte(code), nil)
	if err != nil {
		return nil, err
	}

	return &TokenStore{
		filePath: filePath,
		key:      key,
		salt:     salt,
	}, nil
}

func CreateKey(code []byte, salt []byte) ([]byte, []byte, error) {
	if salt == nil {
		salt = make([]byte, 16)
		_, err := io.ReadFull(rand.Reader, salt)
		if err != nil {
			return nil, nil, errors.New("Could not create token key")
		}
	}

	k := sha256.New()
	k.Write(salt)
	k.Write(code)
	return k.Sum(nil), salt, nil
}

func (store *TokenStore) SaveTokens(tokens *Token, salt []byte) error {
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

	storageData := EncryptStore{
		Data: encrypt,
		IV:   iv,
		Salt: salt,
	}

	jsonData, err := json.Marshal(storageData)
	if err != nil {
		return err
	}

	return os.WriteFile(store.filePath, jsonData, 0600)
}

func (store *TokenStore) LoadTokens() (*Token, error) {
	encrypt, err := os.ReadFile(store.filePath)
	if err != nil {
		return nil, err
	}

	var storage EncryptStore
	err = json.Unmarshal(encrypt, &storage)
	if err != nil {
		return nil, err
	}

	log.Println("Data: ", string(storage.Data))
	log.Println("IV: ", string(storage.IV))
	log.Println("Salt: ", string(storage.Salt))

	// unencrypt data

	return nil, nil
}
