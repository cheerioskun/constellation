package provisioner

import "constellation/internal/models"

// ProvisionRequest is the request object for the provisioner
// It contains the necessary information to image a new node
type ProvisionRequest struct {
	// ID is the ID of the request
	ID string
	// Node is the node to be provisioned
	Node *models.Node
}

// ProvisionResponse is the response object for the provisioner
// It contains the result of the provisioning process
type ProvisionResponse struct {
	// ID is the ID of the request
	ID     string
	NodeID string
	// Error is the error that occurred during the provisioning process
	Error error
	// TimeTaken is the time taken to provision the node in nanoseconds
	TimeTaken int64
}
