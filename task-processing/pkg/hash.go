package pkg

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func Hash(data string, secretKey string) (string, error) {
	hasher := hmac.New(sha256.New, []byte(secretKey))
	_, err := hasher.Write([]byte(data))
	if err != nil {
		return "", err
	}
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha, nil
}

func VerifyHash(workerID, hash, secretKey string) bool {
	serverHash, err := Hash(workerID, secretKey)
	if err != nil {
		return false
	}
	return serverHash == hash
}
