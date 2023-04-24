package yaml

import (
	"io"

	"gopkg.in/yaml.v3"
)

func Marshal(in interface{}) ([]byte, error) {
	return yaml.Marshal(in)
}

func Unmarshal(in []byte, out interface{}) error {
	return yaml.Unmarshal(in, out)
}

func NewDecoder(r io.Reader) *yaml.Decoder {
	return yaml.NewDecoder(r)
}
