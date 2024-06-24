package store

import (
	"encoding/json"
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

func (s *Store) Save() error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	data, err := json.Marshal(s.nodes)
	if err != nil {
		return err
	}

	return os.WriteFile(s.stateFile, data, 0644)
}

func (s *Store) GetNode(id string) (*models.Node, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	node, ok := s.nodes[id]
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

	s.nodes[node.ID] = node
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

	var nodes []*models.Node
	for _, node := range s.nodes {
		if node.Status == models.StatusUninitialized {
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func (s *Store) GetInitializedNodes() []*models.Node {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var nodes []*models.Node
	for _, node := range s.nodes {
		if node.Status == models.StatusInitialized {
			nodes = append(nodes, node)
		}
	}
	return nodes
}
