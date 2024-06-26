package provisioner

import (
	"fmt"
	"log"
	"time"

	"constellation/internal/models"
	"constellation/internal/monitor"
	"constellation/internal/store"
	"constellation/pkg/ipmi"
)

type Provisioner struct {
	store   *store.Store
	isoPath string
}

func New(store *store.Store, isoPath string) *Provisioner {
	return &Provisioner{
		store:   store,
		isoPath: isoPath,
	}
}

func (p *Provisioner) Provision(node *models.Node) error {
	log.Printf("Starting provisioning for node %s", node.Hostname)
	node.Status = models.StatusProvisioning
	p.store.UpdateNode(node)

	log.Printf("Creating ipmi client to spawn shell")
	ipmiClient := ipmi.NewClient()
	if err := ipmiClient.Connect(node.IPMIIP, node.IPMICreds.Username, node.IPMICreds.Password); err != nil {
		err = fmt.Errorf("failed to connect to IPMI: %v", err)
		return err
	}
	defer ipmiClient.Disconnect()

	log.Printf("IPMI Shell Connected! Mounting ISO on node %s", node.Hostname)
	if err := ipmiClient.MountISO(p.isoPath); err != nil {
		err = fmt.Errorf("failed to mount ISO: %v", err)
		log.Printf("%v", err)
		return err
	}

	log.Printf("ISO Mounted! Power cycling node %s", node.Hostname)
	if err := ipmiClient.PowerCycle(); err != nil {
		err = fmt.Errorf("failed to power cycle node: %v", err)
		return err
	}

	// Wait for the node to ping back
	for i := 0; i < 30; i++ {
		time.Sleep(1 * time.Minute)
		log.Printf("Checking if node %s is reachable", node.Hostname)
		monitor.CheckNodeHealth(node, p.store)
		if node.Status == models.StatusInitialized {
			log.Printf("Node %s is reachable", node.Hostname)
			return nil
		}
	}

	node.Status = models.StatusUnhealthy
	p.store.UpdateNode(node)
	return fmt.Errorf("provisioning timed out for node %s", node.Hostname)
}
