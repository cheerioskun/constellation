# Constellation
Constellation is an automated system for provisioning, monitoring, and managing datacenter nodes with failed or end-of-life SATADOM drives, utilizing IPMI for remote management and custom ISO deployment.
## Overview
Constellation provides a suite of tools to streamline the process of harnessing unused datacenter nodes. It automates the provisioning process, monitors node health, and offers both a CLI and API for easy management.
## Key Components

- Provisioner: Handles the automated provisioning of nodes using IPMI and custom ISO deployment.
- Monitor: Continuously checks the health of provisioned nodes and updates their status.
- API: Provides RESTful endpoints for node management and status updates.
- CLI: Offers command-line interface for easy interaction with the system.
- IPMI Client: Manages low-level interaction with node hardware through IPMI.

## Features

- Automated node provisioning using custom ISO
- Real-time node health monitoring
- RESTful API for integration with other systems
- Command-line interface for manual management
- Configurable using TOML for easy setup and modification

## Getting Started

- Clone the repository
- Create a config.yaml
- Configure your nodes in config.yaml
- Build the project
- Run the Constellation server in bg
- Use stars cli for interacting

## Configuration
Constellation uses a YAML configuration file to manage settings and node information. See config.yaml.example for a sample configuration.
## Usage
### CLI
Use the CLI for quick management tasks:
```
stars list
stars provision <node-id>
```
## API
The API provides endpoints for programmatic interaction:

- GET /nodes: List all nodes
- GET /node?id=<node-id>: Get specific node details
- POST /provision: Initiate node provisioning