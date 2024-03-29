// Copyright (c) 2020 Wind River Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software  distributed
// under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
// OR CONDITIONS OF ANY KIND, either express or implied.

// Yaml wrapper package implements ccli yaml data structures and provides standard library functionality
package yaml

import (
	"io"

	"gopkg.in/yaml.v3"
)

// Marshal() gives the yaml encoding of the data
func Marshal(in interface{}) ([]byte, error) {
	return yaml.Marshal(in)
}

// Unmarshal() parses the yaml encoded data
func Unmarshal(in []byte, out interface{}) error {
	return yaml.Unmarshal(in, out)
}

// NewDecoder() gives a new decoder which read from the io.reader
func NewDecoder(r io.Reader) *yaml.Decoder {
	return yaml.NewDecoder(r)
}
