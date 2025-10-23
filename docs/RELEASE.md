# Release Process

This document describes the automated release process for mcp-vosdroits.

## Overview

The release workflow is triggered automatically when a new tag matching the pattern `v*` is pushed to the repository (e.g., `v1.0.0`, `v2.1.3`).

## Workflow Steps

### 1. Build Multi-Platform Binaries

The workflow builds binaries for the following platforms:

- **Linux**: amd64, arm64
- **macOS**: amd64 (Intel), arm64 (Apple Silicon)
- **Windows**: amd64

All binaries are:
- Statically linked (CGO_ENABLED=0)
- Stripped of debug information (-w -s)
- Built with trimmed paths for reproducibility
- Versioned using the Git tag

### 2. Generate Attestations

Each binary receives a build provenance attestation using GitHub's artifact attestation feature. This provides:
- Verification of the build environment
- Cryptographic proof of artifact authenticity
- Supply chain security compliance

### 3. Create GitHub Release

The workflow automatically:
- Downloads all built binaries
- Generates SHA256 checksums for all artifacts
- Creates a GitHub Release with auto-generated release notes
- Attaches all binaries and checksums as release assets

### 4. Build and Push Docker Images

The Docker build process:
- Reuses the pre-built Linux binaries (amd64 and arm64)
- Avoids rebuilding the Go application
- Builds multi-platform Docker images
- Pushes to GitHub Container Registry (ghcr.io)

Docker images are tagged with:
- Full semantic version (e.g., `1.0.0`)
- Major.minor version (e.g., `1.0`)
- Major version (e.g., `1`)
- `latest` tag

## Creating a Release

### 1. Version the Release

Update version numbers if needed in:
- `README.md`
- Any other version references

### 2. Create and Push the Tag

```bash
# Create a new tag
git tag v1.0.0

# Push the tag to GitHub
git push origin v1.0.0
```

### 3. Monitor the Workflow

1. Go to the Actions tab in GitHub
2. Watch the "Release" workflow execution
3. Verify all jobs complete successfully

### 4. Verify the Release

After the workflow completes:

1. **Check GitHub Releases**: Navigate to the Releases page and verify:
   - Release notes are generated
   - All binaries are attached
   - Checksums file is present

2. **Verify Docker Images**: Check the Packages section:
   - All tags are created
   - Images are available for both amd64 and arm64

3. **Test a Binary**:
   ```bash
   # Download a binary
   curl -LO https://github.com/guigui42/mcp-vosdroits/releases/download/v1.0.0/mcp-vosdroits-linux-amd64
   
   # Verify checksum
   sha256sum mcp-vosdroits-linux-amd64
   
   # Compare with checksums.txt from release
   curl -LO https://github.com/guigui42/mcp-vosdroits/releases/download/v1.0.0/checksums.txt
   grep linux-amd64 checksums.txt
   ```

4. **Test Docker Image**:
   ```bash
   docker pull ghcr.io/guigui42/mcp-vosdroits:1.0.0
   docker run -i ghcr.io/guigui42/mcp-vosdroits:1.0.0
   ```

## Verify Artifact Attestations

You can verify the authenticity of released binaries using the GitHub CLI:

```bash
# Install GitHub CLI if needed
# https://cli.github.com/

# Verify a binary
gh attestation verify mcp-vosdroits-linux-amd64 \
  -R guigui42/mcp-vosdroits
```

## Troubleshooting

### Release Workflow Fails

1. Check the workflow logs in the Actions tab
2. Common issues:
   - Build failures: Check Go version compatibility
   - Test failures: Ensure all tests pass locally
   - Docker build issues: Verify Dockerfile.release syntax

### Missing Binaries

If a platform binary is missing:
1. Check if the build matrix includes that platform
2. Verify the build step completed for that platform
3. Check artifact upload logs

### Docker Image Issues

If Docker images aren't published:
1. Verify GitHub token has `packages: write` permission
2. Check Docker login step succeeded
3. Verify the binary artifacts were downloaded correctly

## Release Checklist

Before creating a release:

- [ ] All tests pass locally (`make test`)
- [ ] Code is formatted (`make fmt`)
- [ ] Static analysis passes (`make vet`)
- [ ] Documentation is updated
- [ ] CHANGELOG is updated (if maintained)
- [ ] Version numbers are updated where needed
- [ ] Tag follows semantic versioning (vMAJOR.MINOR.PATCH)

## Rollback

If a release needs to be rolled back:

1. Delete the GitHub Release (this keeps the tag)
2. Delete the Docker images from the Packages section
3. Create a new patch release with the fix
4. Optionally delete the tag:
   ```bash
   git push --delete origin v1.0.0
   git tag -d v1.0.0
   ```

## Security Considerations

- All binaries are built in GitHub Actions with a clean environment
- Attestations provide cryptographic proof of build provenance
- Docker images run as non-root user
- No secrets or sensitive data are included in releases
- All artifacts are signed and verifiable

## Automation Benefits

This release process provides:

1. **Consistency**: Same build process every time
2. **Security**: Attestations and reproducible builds
3. **Efficiency**: No manual binary building or uploading
4. **Multi-platform**: Automatic cross-compilation
5. **Docker Integration**: Reuses binaries to avoid rebuild
6. **Transparency**: Full audit trail in GitHub Actions
