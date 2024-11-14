# Release Process

This document describes the release process for Inbox451.

## Version Scheme

We follow [Semantic Versioning](https://semver.org/):
- MAJOR version for incompatible API changes
- MINOR version for new functionality in a backward compatible manner
- PATCH version for backward compatible bug fixes

## Release Checklist

### 1. Preparation

1. Ensure all tests pass (TODO: Add tests)
   ```bash
   make test
   make test-integration
   ```

### 2. Create Release

1. Create and push a new tag:
   ```bash
   make release-tag VERSION=v1.0.0 # Replace with new version
   ```

2. GoReleaser will automatically:
   - Build binaries for all platforms
   - Create Docker images
   - Generate release notes
   - Upload artifacts to GitHub

### 3. Verify Release

1. Check GitHub release page
2. Verify Docker images:
   ```bash
   docker pull ghcr.io/inbox451/inbox451:latest
   docker pull ghcr.io/inbox451/inbox451:v1.0.0
   ```

3. Test binary downloads:
   ```bash
   # Linux amd64
   curl -L https://github.com/inbox451/inbox451/releases/download/v1.0.0/inbox451_Linux_x86_64.tar.gz | tar xz

   # macOS arm64
   curl -L https://github.com/inbox451/inbox451/releases/download/v1.0.0/inbox451_Darwin_arm64.tar.gz | tar xz
   ```

## Available Artifacts

Each release provides:

### Binaries
- Linux (amd64, arm64)
- macOS (amd64, arm64)

### Docker Images
- `ghcr.io/inbox451/inbox451:v1.0.0`
- `ghcr.io/inbox451/inbox451:latest`

### Checksums
- SHA256 checksums for all artifacts

## Hotfix Process

For urgent fixes:

1. Create hotfix branch from the release tag:
   ```bash
   git checkout -b hotfix/v1.0.1 v1.0.0
   ```

2. Make necessary fixes

3. Follow normal release process with new patch version

## Release Automation

Releases are automated using:
- GitHub Actions for CI/CD
- GoReleaser for build and publish
- Docker Buildx for multi-arch images

See `.github/workflows/release.yml` and `.goreleaser.yaml` for implementation details.
