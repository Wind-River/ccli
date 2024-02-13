// This package defines data structure for expected configuration files
package config

type ConfigData struct {
	ServerAddr string `mapstructure:"server_addr"`
	ServerPort int64  `mapstructure:"server_port"`
	LogFile    string `mapstructure:"log_file"`
	LogLevel   int64  `mapstructure:"log_level"`
	JsonIndent int64  `mapstructure:"json_indent"`
}
