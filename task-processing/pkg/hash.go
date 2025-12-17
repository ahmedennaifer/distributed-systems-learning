package pkg

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

type RegisterPayload struct {
	Hash     string
	WorkerID string
}

func Hash(data string, secretKey string) (string, error) {
	hasher := hmac.New(sha256.New, []byte(secretKey))
	_, err := hasher.Write([]byte(data))
	if err != nil {
		return "", err
	}
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha, nil
}

func VerifyHash(payload RegisterPayload, secretKey string) bool {
	serverHash, err := Hash(payload.WorkerID, secretKey)
	if err != nil {
		return false
	}
	return serverHash == payload.Hash
}
