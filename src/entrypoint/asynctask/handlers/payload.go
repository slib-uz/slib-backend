package handlers

import (
	"encoding/json"
	"fmt"

	"slib.uz/src/core/domain/entity"
)

func UnmarshalPayload[T any](task *entity.AsyncTask) (T, error) {
	var wrapper struct {
		Payload T `json:"payload"`
	}
	if err := json.Unmarshal(task.Payload, &wrapper); err != nil {
		var zero T
		return zero, fmt.Errorf("unmarshal task payload: %w", err)
	}
	return wrapper.Payload, nil
}
