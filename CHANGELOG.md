# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.1.0] - 2025-01-02

### Added
- Infrastructure migration capabilities
  - New `migrate` command for handling infrastructure changes
  - Support for migrating nodes between roles
  - Traefik migration between nodes
  - Manager node setup from existing infrastructure
- Role-based server configuration
  - Dedicated roles for Manager, GitLab, Monitor, and Apps nodes
  - Automatic role-specific setup and configuration
  - Role-based labels and constraints
- Zero-downtime Traefik migration
  - Configuration preservation during migration
  - SSL certificate migration
  - Automatic rollback on failure
  - Service routing continuity

### Changed
- Enhanced server role management
  - Improved role validation
  - Better error handling for role changes
  - More detailed logging during migrations
- Updated deployment process
  - Better handling of existing configurations
  - Improved backup procedures
  - More robust error recovery

### Fixed
- Issue with variable redeclaration in command files
- Compilation errors related to undefined variables
- Error handling in Execute function

## [1.0.0] - 2025-01-02

### Added
- Initial release of the Infrastructure Management CLI
- Authentication system with secure credential storage using SQLite and AES-GCM encryption
- Swarm management commands (init, join, stats)
- Service management with support for different application types
- Monitoring capabilities with health checks and service status
- Backup management system
- Global configuration system with YAML support
- Comprehensive documentation and examples

### Changed
- Unified command structure with consistent flag naming
- Improved error handling and user feedback
- Standardized logging format

### Fixed
- Variable declaration conflicts in command files
- Build system issues with undefined variables
- Execute function return type in root command

## [0.1.0] - 2024-12-15

### Added
- Basic project structure
- Command-line interface framework using Cobra
- Initial implementation of core commands
- Basic documentation
