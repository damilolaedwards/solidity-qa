package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func CreateRequest(url string, requestBody any, method string, headers map[string]string) (*http.Request, error) {
	requestBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Add(k, v)
	}
	return req, nil
}

func DoRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// FetchFileContent fetches the content of a file from a given URL and returns it as a string.
func FetchFileContent(url string, headers map[string]string) (string, error) {
	// Create a new GET request using the existing CreateRequest function.
	req, err := CreateRequest(url, nil, http.MethodGet, headers)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Execute the request using the existing DoRequest function.
	resp, err := DoRequest(req)
	if err != nil {
		return "", fmt.Errorf("failed to perform request: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error closing response body: ", err)
		}
	}(resp.Body)

	// Check if the status code indicates a successful response.
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	// Read the response body.
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Convert the body to a string and return it.
	return string(bodyBytes), nil
}
