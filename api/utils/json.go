package utils

import (
	"encoding/json"
	"net/http"
)

func DecodeRequestBody(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return err
	}
	return nil
}
