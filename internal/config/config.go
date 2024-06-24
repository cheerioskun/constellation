package config

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	StateFilePath   string `toml:"state_file_path"`
	ISOPath         string `toml:"iso_path"`
	IPMIToolJarPath string `toml:"ipmi_tool_jar_path"`
	ListenAddress   string `toml:"listen_address"`
	Nodes           []Node `toml:"nodes"`
}

type Node struct {
	IPMIIP    string    `toml:"ipmi_ip"`
	IPMICreds IPMICreds `toml:"ipmi_creds"`
	NodeIP    string    `toml:"node_ip"`
}

type IPMICreds struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
}

func Load(filename string) (*Config, error) {
	var cfg Config
	if _, err := toml.DecodeFile(filename, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
