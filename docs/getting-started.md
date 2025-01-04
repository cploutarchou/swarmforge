# Getting Started with SwarmForge

## Prerequisites

- Docker installed and running
- Go 1.21 or later
- Access to target servers for swarm deployment

## Installation

Install SwarmForge using Go:

```bash
go install github.com/cploutarchou/swarmforge@latest
```

## Basic Configuration

1. **Environment Setup**
   - Ensure Docker daemon is running
   - Configure SSH access to target servers
   - Set up necessary firewall rules

2. **Initial Configuration**
   - Create a configuration file (optional)
   - Set up environment variables (if needed)

## First Steps

### 1. Initialize a New Swarm

```bash
infra swarm init --ip <manager-ip> --role manager
```

This command:
- Initializes a new Docker Swarm
- Sets up the first manager node
- Configures initial networking

### 2. Configure Traefik

```bash
infra setup traefik --ip <manager-ip> --domain example.com --email admin@example.com
```

This sets up:
- Traefik as reverse proxy
- Automatic SSL certificate management
- Basic routing configuration

### 3. Add Additional Nodes

```bash
infra swarm join --ip <node-ip> --manager-ip <manager-ip> --role <gitlab|monitor|apps>
```

Available roles:
- `gitlab`: GitLab server node
- `monitor`: Monitoring server node
- `apps`: Application deployment node

## Verification

After setup, verify your installation:

1. Check node status
2. Verify Traefik configuration
3. Test basic connectivity
4. Validate SSL certificates
