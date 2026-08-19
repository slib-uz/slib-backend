package network

import (
	"encoding/json"
	"net/http"

	"slib.uz/src/core/utils"
)

func GetBody[T any](resp *http.Response) (*T, error) {
	var body T
	defer utils.Closer(resp.Body)()
	err := json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, err
	}
	return &body, nil
}

func GetErrorBody(resp *http.Response) map[string]any {
	var body map[string]any
	defer utils.Closer(resp.Body)()
	_ = json.NewDecoder(resp.Body).Decode(&body)
	return body
}
