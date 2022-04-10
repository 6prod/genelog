// Package JSON formats log outputs into JSON using
// the structure { "context": {}, "message": "message" }
//
// Note that the encoding/json package used underneath
// does not work with embbeded structures until the go
// issue https://github.com/golang/go/issues/6213
// is solved.
//
// A temporary solution is to define the MarshalJSON() ([]byte, error)
// method on the context type.
//
// Example:
//
// func (c Context) MarshalJSON() ([]byte, error) {
// 	return libjson.Marshal(struct {
// 		Time  libtime.Time `json:"time"`
// 		Level level.Level  `json:"level"`
// 	}{
// 		Time:  c.WithTime.Time(),
// 		Level: c.WithLevel.Level(),
// 	})
// }
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
