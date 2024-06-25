package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"constellation/internal/config"
	"constellation/internal/store"
)

var rootCmd = &cobra.Command{
	Use:   "constellation",
	Short: "Constellation is a tool for managing datacenter nodes",
	Long:  `Constellation allows you to provision, monitor, and manage nodes in your datacenter.`,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all nodes in the constellation",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load("config.yaml")
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}

		store := store.New(cfg.StateFilePath)
		if err := store.Load(); err != nil {
			fmt.Printf("Error loading state: %v\n", err)
			os.Exit(1)
		}

		for _, node := range store.GetAllNodes() {
			fmt.Printf("ID: %s, IPMI IP: %s, Node IP: %s, Status: %d\n", node.Hostname, node.IPMIIP, node.NodeIP, node.Status)
		}
	},
}

var provisionCmd = &cobra.Command{
	Use:   "provision [nodeID]",
	Short: "Provision a specific node",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		nodeID := args[0]
		fmt.Printf("Provisioning node %s...\n", nodeID)
		// Implement provisioning logic here
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(provisionCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
