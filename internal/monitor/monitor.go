package monitor

import (
	"log"
	"net"
	"time"

	"constellation/internal/models"
	"constellation/internal/store"
)

type Monitor struct {
	store *store.Store
}

func New(store *store.Store) *Monitor {
	return &Monitor{
		store: store,
	}
}

func (m *Monitor) MonitorAll() {
	for {
		nodes := m.store.GetInitializedNodes()
		for _, node := range nodes {
			go CheckNodeHealth(node, m.store)
		}
		time.Sleep(5 * time.Minute)
	}
}

func (m *Monitor) Monitor(nodeName string) {
	for {
		node, ok := m.store.GetNode(nodeName)
		if ok {
			CheckNodeHealth(node, m.store)
			m.store.UpdateNode(node)
		}
		time.Sleep(5 * time.Minute)
	}
}

func CheckNodeHealth(node *models.Node, store *store.Store) {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(node.NodeIP, "22"), 5*time.Second)
	if err != nil {
		log.Printf("Node %s is not reachable: %v", node.Hostname, err)
		node.Status = models.StatusUnhealthy
	} else {
		conn.Close()
		node.Status = models.StatusInitialized
	}
	node.LastCheck = time.Now()
	store.UpdateNode(node)
}
