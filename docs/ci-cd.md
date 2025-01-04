# CI/CD Pipeline Documentation

## Overview

The CI/CD pipeline automates testing, building, deployment, and release management of SwarmForge, including GitHub synchronization and Go package publishing.

## Pipeline Stages

### 1. Test Stage
- Runs unit tests
- Performs code linting
- Executes integration tests
- Generates test coverage reports

### 2. Build Stage
- Compiles for multiple platforms
- Generates binaries
- Creates artifacts
- Performs dependency checks

### 3. Package Stage
- Creates distribution packages
- Bundles documentation
- Prepares release artifacts
- Generates checksums

### 4. Publish Stage
- Publishes to pkg.go.dev
- Updates Go package index
- Verifies package availability

### 5. Sync Stage
- Synchronizes with GitHub repository
- Mirrors all commits
- Syncs tags and branches

### 6. Release Stage
- Creates GitHub releases
- Uploads release assets
- Updates release notes
- Manages versioning

## Version Management

### Semantic Versioning
We follow semantic versioning (SemVer) for releases:
- MAJOR version for incompatible API changes
- MINOR version for backward-compatible functionality
- PATCH version for backward-compatible bug fixes

Example: v1.2.3
- 1: Major version
- 2: Minor version
- 3: Patch version

### Release Process

1. Update version in project files
2. Update CHANGELOG.md with changes
3. Create and push a new tag:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

The pipeline automatically:
- Runs all tests
- Builds binaries
- Creates packages
- Publishes to pkg.go.dev
- Syncs to GitHub
- Creates GitHub release

## Go Package Publishing

The pipeline automatically publishes to pkg.go.dev when a new version tag is pushed:

1. Tag format must be `vX.Y.Z` (e.g., v1.0.0)
2. Package is published to proxy.golang.org
3. pkg.go.dev automatically indexes the new version

## GitHub Integration

### Synchronization
- All commits are mirrored to GitHub
- Tags are synchronized
- Branches are maintained
- Release notes are generated from CHANGELOG.md

### Release Creation
- Automated GitHub release creation
- Release assets are uploaded
- Release notes from CHANGELOG.md
- Version tags are properly managed

## Configuration Details

The pipeline is configured in `.gitlab-ci.yml`:

```yaml
stages:
  - test
  - build
  - package
  - publish
  - sync
  - release

variables:
  GO_VERSION: "1.21"
  CGO_ENABLED: "0"
  GITHUB_REPO: "github.com/cploutarchou/swarmforge"
  GITLAB_REPO: "gitlab.com/cploutarchou/swarmforge"
```

## Required Variables

Configure these variables in GitLab CI/CD Settings:

1. `GITHUB_TOKEN`: GitHub personal access token with repo scope
   - Required for GitHub synchronization and release creation
   - Must have `repo` scope permissions
   - Should be protected and masked

## Troubleshooting

Common issues and solutions:

1. **Failed Go Package Publishing**
   - Verify tag format (must be vX.Y.Z)
   - Check go.mod configuration
   - Ensure GOPROXY is accessible

2. **GitHub Sync Issues**
   - Verify GITHUB_TOKEN permissions
   - Check repository access
   - Validate remote URLs

3. **Release Creation Problems**
   - Verify CHANGELOG.md format
   - Check tag format
   - Review GitHub API responses
