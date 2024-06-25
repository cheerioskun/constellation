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
	Hostname  string     `json:"hostname" yaml:"hostname"`
	IPMIIP    string     `json:"ipmi_ip" yaml:"ipmi_ip"`
	NodeIP    string     `json:"node_ip" yaml:"node_ip"`
	Status    NodeStatus `json:"status" yaml:"status"`
	LastCheck time.Time  `json:"last_check" yaml:"last_check"`
	CPU       int        `json:"cpu" yaml:"cpu"`
	Memory    int        `json:"memory" yaml:"memory"`
	IPMICreds IPMICreds  `json:"ipmi_creds" yaml:"ipmi_creds"`
}

type IPMICreds struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
