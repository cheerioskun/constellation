package models

import "time"

type NodeStatus string

const (
	StatusUninitialized NodeStatus = "uninitialized"
	StatusProvisioning  NodeStatus = "provisioning"
	StatusInitialized   NodeStatus = "initialized"
	StatusUnhealthy     NodeStatus = "unhealthy"
)

type Node struct {
	ID        string     `json:"id"`
	IPMIIP    string     `json:"ipmi_ip"`
	NodeIP    string     `json:"node_ip"`
	Status    NodeStatus `json:"status"`
	LastCheck time.Time  `json:"last_check"`
	CPU       int        `json:"cpu"`
	Memory    int        `json:"memory"`
	Hostname  string     `json:"hostname"`
	IPMICreds IPMICreds  `json:"-"`
}

type IPMICreds struct {
	Username string
	Password string
}
