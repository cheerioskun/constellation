package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/nodes", s.ListNodes).Methods("GET")
	r.HandleFunc("/node", s.GetNode).Methods("GET")
	r.HandleFunc("/provision", s.ProvisionNode).Methods("POST")
	r.HandleFunc("/pingback", s.PingBack).Methods("POST")

	return r
}

func (s *Server) Start(addr string) error {
	r := s.SetupRoutes()
	return http.ListenAndServe(addr, r)
}

func (s *Server) Shutdown() error {
	// Nothing here right now
	return nil
}
