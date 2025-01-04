# Infrastructure Setup CLI

A powerful CLI tool for managing Docker Swarm infrastructure, handling deployments, and managing services.

## Features

- **Infrastructure Management**

  - Swarm initialization and node management
  - Role-based node configuration (Manager, GitLab, Monitor, Apps)
  - Secure credential storage and management
  - Infrastructure migration and upgrades

- **Service Management**

  - Automated service deployment
  - Configuration management
  - Service health monitoring
  - Backup and restore capabilities

- **Traefik Integration**
  - Automatic SSL certificate management
  - Dynamic routing configuration
  - Multi-domain support
  - Zero-downtime migrations

## Documentation

Comprehensive documentation is available in the [docs](docs/) directory:

- [Getting Started Guide](docs/getting-started.md)
- [Architecture Overview](docs/architecture.md)
- [Configuration Guide](docs/configuration.md)

## Installation

```bash
go install github.com/cploutarchou/swarmforge@latest
```

## Quick Start

1. Initialize a new swarm:

```bash
infra swarm init --ip <manager-ip> --role manager
```

2. Set up Traefik:

```bash
infra setup traefik --ip <manager-ip> --domain example.com --email admin@example.com
```

3. Join additional nodes:

```bash
infra swarm join --ip <node-ip> --manager-ip <manager-ip> --role <gitlab|monitor|apps>
```

## Infrastructure Migration

### Migrating Nodes

Convert existing nodes to new roles while preserving data:

```bash
# Migrate a node to a new role
infra migrate node --ip <node-ip> --role <new-role>

# Migrate with force option (skip confirmation)
infra migrate node --ip <node-ip> --role <new-role> --force
```

### Setting Up Manager Node

Convert an existing node to manager or set up a new manager node:

```bash
# Basic manager setup
infra migrate setup-manager --target-ip <new-manager-ip>

# Setup manager and migrate Traefik
infra migrate setup-manager \
  --target-ip <new-manager-ip> \
  --source-ip <old-traefik-ip> \
  --migrate-traefik
```

### Traefik Migration

Migrate Traefik between nodes while preserving configurations:

```bash
# Migrate Traefik to new node
infra migrate traefik \
  --source-ip <current-traefik-ip> \
  --target-ip <new-traefik-ip>
```

The migration process includes:

- Configuration backup
- SSL certificate preservation
- Zero-downtime transition
- Automatic rollback on failure

## Server Roles

### Manager Node

- Swarm management
- Traefik routing
- SSL certificate management

### GitLab Node

- GitLab server
- CI/CD pipelines
- Repository management

### Monitor Node

- Prometheus metrics
- Grafana dashboards
- Log aggregation

### Apps Node

- Application workloads
- Service deployments
- Data persistence

## Security Features

- Access control
- Secure communication
- SSL/TLS encryption
- Role-based access control

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
