// Json wrapper package implements ccli json data structures and provides standard library functionality
package json

import (
	"encoding/json"
)

type RawMessage json.RawMessage

func Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func MarshalIndent(v any, prefix string, indent string) ([]byte, error) {
	return json.MarshalIndent(v, prefix, indent)
}
