package mapper

import (
	"encoding/json"

	"gorm.io/datatypes"
)

func FromGormJson[T any](in datatypes.JSON) T {
	var out T

	if in == nil {
		return out
	}

	// First, try direct unmarshal to target type
	err := json.Unmarshal(in, &out)
	if err == nil {
		return out
	}

	// If that fails, the data might be stored as a JSON-escaped string (e.g., "{\"key\": \"value\"}")
	// Try to unmarshal as string first, then unmarshal the result
	var jsonStr string
	if err2 := json.Unmarshal(in, &jsonStr); err2 == nil {
		// We got a string, now try to parse it as JSON
		if err3 := json.Unmarshal([]byte(jsonStr), &out); err3 == nil {
			return out
		}
	}

	// If both attempts fail, return zero value instead of panicking
	// This allows the application to continue even if one field has invalid JSON
	return out
}

func ToGormJson(in any) datatypes.JSON {
	_json, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
	return _json
}
