package awshelpers

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

var keyArn = "arn:aws:kms:us-east-1:957015889457:key/d7e1ec89-9186-48b9-98d9-c7d9b9d216d2"
var awsCfg aws.Config

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("failed to load aws config")
	}

	awsCfg = cfg
}

func EncryptKey(token string) (string, error) {
	client := kms.NewFromConfig(awsCfg)

	encryptInput := &kms.EncryptInput{
		KeyId:     aws.String(keyArn),
		Plaintext: []byte(token),
	}

	result, err := client.Encrypt(context.TODO(), encryptInput)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt: %v", err)
	}

	return base64.StdEncoding.EncodeToString(result.CiphertextBlob), nil
}

func DecryptKey(encryptedToken string) (string, error) {
	client := kms.NewFromConfig(awsCfg)

	token, err := base64.StdEncoding.DecodeString(encryptedToken)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt: %v", err)
	}

	decryptInput := &kms.DecryptInput{
		CiphertextBlob: token,
	}

	result, err := client.Decrypt(context.TODO(), decryptInput)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt: %v", err)
	}

	return string(result.Plaintext), nil
}

