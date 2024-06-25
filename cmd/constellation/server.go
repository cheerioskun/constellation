package main

import (
	"constellation/internal/api"
	"constellation/internal/config"
	"constellation/internal/models"
	"constellation/internal/monitor"
	"constellation/internal/provisioner"
	"constellation/internal/store"
	"constellation/pkg/ipmi"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alitto/pond"
	"github.com/davecgh/go-spew/spew"
)

const (
	Concurrency = 3
)

func main() {

	cfg, err := config.Load("config.yaml")
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
	log.Println(spew.Sdump(cfg.Nodes))
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
	go StartAutoProvisioning(store, prov)

	// Start monitoring initialized nodes
	for _, node := range store.GetInitializedNodes() {
		log.Printf("Monitoring node %s\n", node.Hostname)
		go mon.Monitor(node.Hostname)
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

func StartAutoProvisioning(store *store.Store, prov *provisioner.Provisioner) {

	pool := pond.New(Concurrency, Concurrency*10)
	nodes_to_provision := store.GetUninitializedNodes()
	in_progress_nodes := store.GetNodesByStatus(models.StatusProvisioning)
	nodes_to_provision = append(nodes_to_provision, in_progress_nodes...)
	for _, _node := range nodes_to_provision {
		node := _node
		log.Printf("Requesting provisioning for %s\n", node.Hostname)
		pool.Submit(func() {
			prov.Provision(node)
		})
	}
	pool.StopAndWait()
}
