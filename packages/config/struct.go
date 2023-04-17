// This package defines data structure for expected configuration files
package config

type ConfigData struct {
	ServerAddr string `yaml:"server_addr"`
	ServerPort int64  `yaml:"server_port"`
	LogFile    string `yaml:"log_file"`
	LogLevel   int64  `yaml:"log_level"`
	JsonIndent int64  `yaml:"json_indent"`
}
