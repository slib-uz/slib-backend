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

	err := json.Unmarshal(in, &out)
	if err != nil {
		var jsonStr string
		if err := json.Unmarshal(in, &jsonStr); err == nil {
			if err := json.Unmarshal([]byte(jsonStr), &out); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	return out
}

func ToGormJson(in any) datatypes.JSON {
	_json, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
	return _json
}
