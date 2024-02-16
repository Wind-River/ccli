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

// This package defines data structure for expected configuration files
package config

import (
	"os"

	"github.com/pkg/errors"
)

type ConfigData struct {
	ServerAddr string `mapstructure:"server_addr"`
	ServerPort int64  `mapstructure:"server_port"`
	LogFile    string `mapstructure:"log_file"`
	LogLevel   int64  `mapstructure:"log_level"`
	JsonIndent int64  `mapstructure:"json_indent"`
}

type LogWriter struct {
	Stdout *os.File
	File   *os.File
}

// Write implements io.Writer.
func (logWriter *LogWriter) Write(p []byte) (n int, err error) {
	n, err = logWriter.Stdout.Write(p)
	if err != nil {
		return n, errors.Wrapf(err, "Could not write to stdout")
	}
	n, err = logWriter.File.Write(p)
	return n, errors.Wrapf(err, "Could not write to log file")
}
