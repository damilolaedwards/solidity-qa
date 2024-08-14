package client

import (
	"bytes"
	"encoding/json"
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
