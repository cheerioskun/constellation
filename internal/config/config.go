package config

import (
	"constellation/internal/models"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	StateFilePath   string                 `yaml:"state_file_path"`
	ISOPath         string                 `yaml:"iso_path"`
	IPMIToolJarPath string                 `yaml:"ipmi_tool_jar_path"`
	ListenAddress   string                 `yaml:"listen_address"`
	Nodes           map[string]models.Node `yaml:"nodes"`
}

func Load(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
