package store

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"constellation/internal/models"
)

type Store struct {
	nodes     map[string]*models.Node
	mutex     sync.RWMutex
	stateFile string
}

func New(stateFile string) *Store {
	return &Store{
		nodes:     make(map[string]*models.Node),
		stateFile: stateFile,
	}
}

func (s *Store) Load() error {
	data, err := os.ReadFile(s.stateFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	return json.Unmarshal(data, &s.nodes)
}

func (s *Store) LoadNodes(nodes map[string]models.Node) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, _node := range nodes {
		if s.nodes[_node.Hostname] != nil {
			continue
		}
		log.Printf("Loading node: %s\n", _node.Hostname)
		node := _node
		s.nodes[node.Hostname] = &node
	}
}

func (s *Store) Save() error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	data, err := json.MarshalIndent(s.nodes, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.stateFile, data, 0644)
}

func (s *Store) GetNode(name string) (*models.Node, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	node, ok := s.nodes[name]
	return node, ok
}

func (s *Store) GetNodeByHostname(hostname string) (*models.Node, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, node := range s.nodes {
		if node.Hostname == hostname {
			return node, true
		}
	}
	return nil, false
}

func (s *Store) UpdateNode(node *models.Node) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.nodes[node.Hostname] = node
}

func (s *Store) GetAllNodes() []*models.Node {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var nodes []*models.Node
	for _, node := range s.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

func (s *Store) GetUninitializedNodes() []*models.Node {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.getNodesByStatus(models.StatusUninitialized)
}

func (s *Store) GetInitializedNodes() []*models.Node {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.getNodesByStatus(models.StatusInitialized)
}

func (s *Store) GetNodesByStatus(status models.NodeStatus) []*models.Node {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.getNodesByStatus(status)
}

func (s *Store) getNodesByStatus(status models.NodeStatus) []*models.Node {
	var nodes []*models.Node
	for _, node := range s.nodes {
		if node.Status == status {
			nodes = append(nodes, node)
		}
	}
	return nodes
}
