// config/config.go

package config

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"os"
)

var (
	RpcAddr       string
	ChainId       string
	PrivateKey    string
	PrivateAESKey []byte
)

func Init() {
	RpcAddr = "https://gnfd-testnet-fullnode-tendermint-us.bnbchain.org:443"
	ChainId = "greenfield_5600-1"
	PrivateKey = "86c6252d772b7a85fd566e19d1dab0a7f6b246348bc133689633db4c0322cb14"
	PrivateAESKey, _ = base64.StdEncoding.DecodeString("9Z8wR5lNzKABaZa45jSt7h7J59RHbDm9aLbCFQqKInk=")
	//PrivateAESKey, _ = generateAESKey(256)
}

func generateAESKey(bits int) ([]byte, error) {
	keyLength := bits / 8 // 8 bits in a byte
	key := make([]byte, keyLength)

	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func getEncryptionKeyFromEnv() ([]byte, error) {
	encodedKey := os.Getenv("AES_KEY_ENV_VAR")
	if encodedKey == "" {
		return nil, errors.New("encryption key not found in environment variable")
	}

	key, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil, err
	}

	return key, nil
}
