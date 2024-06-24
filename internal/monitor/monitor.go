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
			go m.checkNodeHealth(node)
		}
		time.Sleep(5 * time.Minute)
	}
}

func (m *Monitor) Monitor(nodeID string) {
	for {
		node, ok := m.store.GetNode(nodeID)
		if ok {
			m.checkNodeHealth(node)
		}
		time.Sleep(5 * time.Minute)
	}
}

func (m *Monitor) checkNodeHealth(node *models.Node) {
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(node.NodeIP, "22"), 5*time.Second)
	if err != nil {
		log.Printf("Node %s is unhealthy: %v", node.ID, err)
		node.Status = models.StatusUnhealthy
	} else {
		conn.Close()
		node.Status = models.StatusInitialized
	}
	node.LastCheck = time.Now()
	m.store.UpdateNode(node)
}
