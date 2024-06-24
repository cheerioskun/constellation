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

	store := store.New(cfg.StateFilePath)
	if err := store.Load(); err != nil {
		log.Fatalf("Failed to load state: %v", err)
	}

	ipmi.SetJarPath(cfg.IPMIToolJarPath)

	// Initialize the provisioner and monitor
	prov := provisioner.New(store, cfg.ISOPath)
	mon := monitor.New(store)

	// Start the API server
	server := api.NewServer(store, prov)
	go func() {
		if err := server.Start(cfg.ListenAddress); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Start provisioning uninitialized nodes
	for _, node := range store.GetUninitializedNodes() {
		go prov.Provision(node)
	}

	// Start monitoring initialized nodes
	for _, node := range store.GetInitializedNodes() {
		go mon.Monitor(node.ID)
	}

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
}
