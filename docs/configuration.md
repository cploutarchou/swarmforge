# Configuration Guide

## Configuration Options

### Basic Configuration

```yaml
swarm:
  manager:
    ip: "192.168.1.100"
    port: 2377
  
  network:
    subnet: "10.0.0.0/24"
    driver: "overlay"

  traefik:
    domain: "example.com"
    email: "admin@example.com"
    ssl: true
```

### Node Configuration

```yaml
nodes:
  - role: gitlab
    ip: "192.168.1.101"
    resources:
      cpu: "2"
      memory: "4G"
  
  - role: monitor
    ip: "192.168.1.102"
    resources:
      cpu: "2"
      memory: "4G"
  
  - role: apps
    ip: "192.168.1.103"
    resources:
      cpu: "4"
      memory: "8G"
```

## Environment Variables

Required environment variables:

```bash
# Docker configuration
DOCKER_HOST=tcp://localhost:2375
DOCKER_TLS_VERIFY=1
DOCKER_CERT_PATH=/path/to/certs

# Swarm configuration
SWARM_MANAGER_IP=192.168.1.100
SWARM_MANAGER_PORT=2377

# SSL configuration
SSL_CERT_PATH=/path/to/ssl
SSL_KEY_PATH=/path/to/ssl/key
```

## Security Configuration

### SSL/TLS Settings

```yaml
ssl:
  enabled: true
  provider: "letsencrypt"
  email: "admin@example.com"
  domains:
    - "example.com"
    - "*.example.com"
```

### Access Control

```yaml
rbac:
  enabled: true
  roles:
    - name: admin
      permissions:
        - "*"
    - name: developer
      permissions:
        - "deploy"
        - "scale"
        - "logs"
```

## Monitoring Configuration

```yaml
monitoring:
  prometheus:
    enabled: true
    retention: "15d"
    scrape_interval: "15s"
  
  grafana:
    enabled: true
    admin_user: "admin"
    dashboards_enabled: true
```

## Backup Configuration

```yaml
backup:
  enabled: true
  schedule: "0 0 * * *"
  retention: "7d"
  storage:
    type: "s3"
    bucket: "backups"
    path: "/swarm"
```

## Network Configuration

```yaml
networks:
  overlay:
    name: "swarm-net"
    driver: "overlay"
    subnet: "10.0.0.0/24"
    
  ingress:
    name: "ingress-net"
    driver: "overlay"
    subnet: "10.0.1.0/24"
```

## Service Defaults

```yaml
services:
  defaults:
    replicas: 2
    update_config:
      parallelism: 1
      delay: "10s"
    restart_policy:
      condition: "on-failure"
      max_attempts: 3
```
