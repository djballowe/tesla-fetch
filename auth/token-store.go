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
	"time"
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

	salt := make([]byte, 16)
	_, err = io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, errors.New("Could not create token key")
	}

	key, err := getKey([]byte(code), salt)
	if err != nil {
		return nil, err
	}

	return &TokenStore{
		key:  key,
		salt: salt,
	}, nil
}

func getKey(code []byte, salt []byte) ([]byte, error) {
	k := sha256.New()
	k.Write(salt)
	k.Write(code)
	return k.Sum(nil), nil
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

	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, jsonData, 0600)
}

func (store *TokenStore) LoadTokens(code string) (*Token, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	encrypt, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("error not exist")
			return nil, errors.New("No token stored")
		}
		return nil, err
	}
	if string(encrypt) == "" {
		return nil, errors.New("No token stored")
	}

	var storage EncryptStore
	err = json.Unmarshal(encrypt, &storage)
	if err != nil {
		return nil, err
	}

	store.key, err = getKey([]byte(code), storage.Salt)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(store.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	decrypt, err := gcm.Open(nil, storage.IV, storage.Data, nil)
	if err != nil {
		return nil, err
	}

	var tokens Token
	err = json.Unmarshal(decrypt, &tokens)
	if err != nil {
		return nil, err
	}

	return &tokens, nil
}

func (store *TokenStore) IsExpired(createdAt time.Time, expiresIn int) bool {
	expiration := createdAt.Add(time.Duration(expiresIn) * time.Second)
	return time.Now().After(expiration)
}

func getConfigFilePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	tfetchDir := filepath.Join(configDir, "tfetch")
	err = os.MkdirAll(tfetchDir, 0700)
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(tfetchDir, "tfetch-tokens.dat")

	return filePath, nil
}
