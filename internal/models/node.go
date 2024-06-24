package models

import "time"

type NodeStatus int

const (
	StatusUninitialized NodeStatus = iota
	StatusProvisioning
	StatusInitialized
	StatusUnhealthy
)

type Node struct {
	ID        string     `json:"id" toml:"id"`
	IPMIIP    string     `json:"ipmi_ip" toml:"ipmi_ip"`
	NodeIP    string     `json:"node_ip" toml:"node_ip"`
	Status    NodeStatus `json:"status" toml:"status"`
	LastCheck time.Time  `json:"last_check" toml:"last_check"`
	CPU       int        `json:"cpu" toml:"cpu"`
	Memory    int        `json:"memory" toml:"memory"`
	Hostname  string     `json:"hostname" toml:"hostname"`
	IPMICreds IPMICreds  `json:"ipmi_creds" toml:"ipmi_creds"`
}

type IPMICreds struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
}
