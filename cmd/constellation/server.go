package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"constellation/internal/api"
	"constellation/internal/config"
	"constellation/internal/monitor"
	"constellation/internal/provisioner"
	"constellation/internal/store"
	"constellation/pkg/ipmi"
)

func main() {

	cfg, err := config.Load("config.toml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Println("Configuration loaded successfully")

	store := store.New(cfg.StateFilePath)
	if err := store.Load(); err != nil {
		log.Fatalf("Failed to load state: %v", err)
	}
	log.Println("State loaded successfully")

	// Load the initial node configuration from the configuration file
	store.LoadNodes(cfg.Nodes)
	log.Println("Nodes loaded successfully")

	ipmi.SetJarPath(cfg.IPMIToolJarPath)
	log.Println("IPMI tool jar path set successfully")

	// Initialize the provisioner and monitor
	prov := provisioner.New(store, cfg.ISOPath)
	mon := monitor.New(store)
	log.Println("Provisioner and monitor initialized successfully")

	// Start the API server
	server := api.NewServer(store, prov)
	go func() {
		if err := server.Start(cfg.ListenAddress); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
	log.Println("API server started successfully")

	// Start provisioning uninitialized nodes
	for _, node := range store.GetUninitializedNodes() {
		log.Printf("Provisioning node %s: %s\n", node.ID, node.Hostname)
		go prov.Provision(node)
	}
	log.Println("Provisioning of uninitialized nodes started")

	// Start monitoring initialized nodes
	for _, node := range store.GetInitializedNodes() {
		log.Printf("Monitoring node %s: %s\n", node.ID, node.Hostname)
		go mon.Monitor(node.ID)
	}
	log.Println("Monitoring of initialized nodes started")

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gracefully...")
	if err := server.Shutdown(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	if err := store.Save(); err != nil {
		log.Fatalf("Failed to save state: %v", err)
	}
	log.Println("State saved successfully")
}
