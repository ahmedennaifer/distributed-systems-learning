package pkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func PostRequest(url string, payload any) ([]byte, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, errors.New("cannot marshal json post data")
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
