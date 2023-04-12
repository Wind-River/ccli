package json

import "encoding/json"

type RawMessage json.RawMessage

func Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}
