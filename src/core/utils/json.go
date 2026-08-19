package utils

import "encoding/json"

func JsonUnmarshal[T any](data []byte) (*T, error) {
	var t T
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, err
	}
	return &t, nil

}

func JsonMarshal[T any](t T) []byte {
	data, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return data

}

func JsonStringify[T any](t T) string {
	data, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return string(data)
}
