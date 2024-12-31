package util

import "encoding/json"

func MarshalString(v any) (string, error) {
	bytes, err := json.Marshal(v)
	return string(bytes), err
}
