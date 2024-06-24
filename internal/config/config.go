package config

import (
	"constellation/internal/models"

	"github.com/BurntSushi/toml"
)

type Config struct {
	StateFilePath   string        `toml:"state_file_path"`
	ISOPath         string        `toml:"iso_path"`
	IPMIToolJarPath string        `toml:"ipmi_tool_jar_path"`
	ListenAddress   string        `toml:"listen_address"`
	Nodes           []models.Node `toml:"nodes"`
}

func Load(filename string) (*Config, error) {
	var cfg Config
	if _, err := toml.DecodeFile(filename, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
