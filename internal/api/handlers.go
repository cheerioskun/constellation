package api

import (
	"encoding/json"
	"net/http"

	"constellation/internal/models"
	"constellation/internal/provisioner"
	"constellation/internal/store"
)

type Server struct {
	store       *store.Store
	provisioner *provisioner.Provisioner
}

func NewServer(store *store.Store, provisioner *provisioner.Provisioner) *Server {
	return &Server{
		store:       store,
		provisioner: provisioner,
	}
}

func (s *Server) ListNodes(w http.ResponseWriter, r *http.Request) {
	nodes := s.store.GetAllNodes()
	json.NewEncoder(w).Encode(nodes)
}

func (s *Server) GetNode(w http.ResponseWriter, r *http.Request) {
	nodeID := r.URL.Query().Get("id")
	if nodeID == "" {
		http.Error(w, "Missing node ID", http.StatusBadRequest)
		return
	}

	node, ok := s.store.GetNode(nodeID)
	if !ok {
		http.Error(w, "Node not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(node)
}

func (s *Server) ProvisionNode(w http.ResponseWriter, r *http.Request) {
	var req struct {
		NodeID string `json:"node_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	node, ok := s.store.GetNode(req.NodeID)
	if !ok {
		http.Error(w, "Node not found", http.StatusNotFound)
		return
	}

	go s.provisioner.Provision(node)

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "provisioning started"})
}

func (s *Server) PingBack(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Hostname string `json:"hostname"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	node, ok := s.store.GetNodeByHostname(req.Hostname)
	if !ok {
		http.Error(w, "Node not found", http.StatusNotFound)
		return
	}

	node.Status = models.StatusInitialized
	s.store.UpdateNode(node)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
