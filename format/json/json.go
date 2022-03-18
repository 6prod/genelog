package json

import (
	"bytes"
	"encoding/json"
)

// FormatJSON format into JSON using the structure
// { "context": v, "message": msg }
func JSON(v interface{}, msg string) (string, error) {
	m := struct {
		Context interface{} `json:"context"`
		Message string      `json:"message"`
	}{
		Context: v,
		Message: msg,
	}

	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(m); err != nil {
		return "", err
	}
	return buf.String(), nil
}
