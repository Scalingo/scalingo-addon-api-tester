package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ProvisioningResponse struct {
	ID      string            `json:"id"`
	Config  map[string]string `json:"config"`
	Message string            `json:"message"`
}

type UpdateResponse struct {
	Config  map[string]string `json:"config"`
	Message string            `json:"message"`
}

func doRequest(method string, url string, payload interface{}) (*http.Response, error) {
	payloadJSON, err := json.Marshal(&payload)
	if err != nil {
		return nil, fmt.Errorf("Fail to encode to JSON: %v", err)
	}
	payloadBuffer := bytes.NewBuffer(payloadJSON)

	req, err := http.NewRequest(method, url, payloadBuffer)
	if err != nil {
		return nil, err
	}

	if req.Method != "DELETE" {
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")
	}
	req.SetBasicAuth(manifest.Username, manifest.Password)

	return http.DefaultClient.Do(req)
}
