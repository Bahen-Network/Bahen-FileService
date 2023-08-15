// config/config.go

package config

import (
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
	PrivateAESKey, _ = getEncryptionKeyFromEnv()
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
