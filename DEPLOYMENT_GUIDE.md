# DevPlatform CLI - Deployment Guide

## Overview

This guide explains how binaries are distributed via GitHub Releases using automated CI/CD workflows.

## Deployment Architecture

```
Git Tag (v1.0.0)
    ↓
GitHub Actions Workflow Triggered
    ↓
GoReleaser Builds Binaries
    ↓
    ├─ Linux (amd64, arm64)
    ├─ macOS (amd64, arm64)
    └─ Windows (amd64)
    ↓
Create GitHub Release
    ↓
Upload Artifacts
    ↓
    ├─ Binaries (tar.gz, zip)
    ├─ Checksums (SHA256)
    ├─ Changelog
    └─ Documentation
```

## Automated Release Process

### 1. Trigger Release

The release process is triggered automatically when a version tag is pushed:

```bash
# Create and push a tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

### 2. GitHub Actions Workflow

The `.github/workflows/release.yml` workflow:

1. **Checkout Code**: Fetches the repository with full history
2. **Setup Go**: Installs Go 1.21 with caching
3. **Run Tests**: Executes all tests to ensure quality
4. **Run GoReleaser**: Builds and releases binaries
5. **Upload Artifacts**: Stores build artifacts

### 3. GoReleaser Configuration

The `.goreleaser.yml` configuration defines:

#### Build Targets

| OS      | Architecture | Binary Format |
|---------|--------------|---------------|
| Linux   | amd64        | tar.gz        |
| Linux   | arm64        | tar.gz        |
| macOS   | amd64        | tar.gz        |
| macOS   | arm64        | tar.gz        |
| Windows | amd64        | zip           |

#### Build Flags

```bash
-s -w                                    # Strip debug info
-X main.Version={{.Version}}             # Set version
-X main.GitCommit={{.Commit}}            # Set git commit
-X main.BuildDate={{.Date}}              # Set build date
```

#### Archive Contents

Each release archive includes:
- Binary executable (`devplatform-cli` or `devplatform-cli.exe`)
- README.md
- LICENSE
- Complete documentation (`docs/**/*`)

## Release Assets

### Binary Naming Convention

```
devplatform-cli_<VERSION>_<OS>_<ARCH>.<FORMAT>
```

Examples:
- `devplatform-cli_1.0.0_Linux_x86_64.tar.gz`
- `devplatform-cli_1.0.0_Darwin_arm64.tar.gz`
- `devplatform-cli_1.0.0_Windows_x86_64.zip`

### Checksums

A `checksums.txt` file is generated with SHA256 hashes for all binaries:

```
sha256sum -c checksums.txt
```

## Installation Methods

### 1. Direct Download (GitHub Releases)

**Linux/macOS:**
```bash
# Download the appropriate binary
wget https://github.com/vedantchimote/DevPlatform-CLI/releases/download/v1.0.0/devplatform-cli_1.0.0_Linux_x86_64.tar.gz

# Extract
tar -xzf devplatform-cli_1.0.0_Linux_x86_64.tar.gz

# Install
sudo mv devplatform-cli /usr/local/bin/

# Verify
devplatform-cli version
```

**Windows:**
```powershell
# Download from GitHub Releases
# Extract the zip file
# Move devplatform-cli.exe to a directory in your PATH
# Verify
devplatform-cli.exe version
```

### 2. Homebrew (macOS/Linux)

```bash
# Add tap (once homebrew-tap repository is set up)
brew tap vedantchimote/tap

# Install
brew install devplatform-cli

# Verify
devplatform-cli version
```

### 3. Package Managers

**Debian/Ubuntu (.deb):**
```bash
wget https://github.com/vedantchimote/DevPlatform-CLI/releases/download/v1.0.0/devplatform-cli_1.0.0_amd64.deb
sudo dpkg -i devplatform-cli_1.0.0_amd64.deb
```

**RHEL/CentOS (.rpm):**
```bash
wget https://github.com/vedantchimote/DevPlatform-CLI/releases/download/v1.0.0/devplatform-cli_1.0.0_amd64.rpm
sudo rpm -i devplatform-cli_1.0.0_amd64.rpm
```

## Verification

### 1. Verify Checksum

```bash
# Download checksums
wget https://github.com/vedantchimote/DevPlatform-CLI/releases/download/v1.0.0/checksums.txt

# Verify (Linux/macOS)
sha256sum -c checksums.txt

# Verify (Windows)
certutil -hashfile devplatform-cli_1.0.0_Windows_x86_64.zip SHA256
```

### 2. Verify Binary

```bash
# Check version
devplatform-cli version

# Check dependencies
devplatform-cli version --check-deps
```

## Release Checklist

Before creating a release:

- [ ] All tests passing
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] Version bumped in code
- [ ] Release notes prepared
- [ ] Manual testing completed

## Monitoring Release

### Check Workflow Status

1. Go to: https://github.com/vedantchimote/DevPlatform-CLI/actions
2. Find the "Release" workflow
3. Monitor the build progress

### Verify Release

1. Go to: https://github.com/vedantchimote/DevPlatform-CLI/releases
2. Verify the release is published
3. Check all assets are present:
   - [ ] Linux amd64 binary
   - [ ] Linux arm64 binary
   - [ ] macOS amd64 binary
   - [ ] macOS arm64 binary
   - [ ] Windows amd64 binary
   - [ ] checksums.txt
   - [ ] .deb package
   - [ ] .rpm package

## Troubleshooting

### Workflow Fails

**Check logs:**
```bash
# View workflow logs on GitHub Actions page
# Common issues:
# - Test failures
# - Build errors
# - Permission issues
```

**Fix and re-release:**
```bash
# Delete the tag locally and remotely
git tag -d v1.0.0
git push origin :refs/tags/v1.0.0

# Fix the issue
# Create tag again
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

### Missing Assets

If some assets are missing:
1. Check GoReleaser logs
2. Verify `.goreleaser.yml` configuration
3. Ensure all build targets are specified

### Permission Issues

Ensure the workflow has proper permissions:
```yaml
permissions:
  contents: write  # Required for creating releases
```

## Post-Release Tasks

After a successful release:

1. **Announce Release**
   - Update documentation site
   - Post on social media
   - Notify users via email/Slack

2. **Update Documentation**
   - Update installation instructions
   - Update version references
   - Update screenshots if needed

3. **Monitor Issues**
   - Watch for bug reports
   - Monitor download statistics
   - Gather user feedback

4. **Plan Next Release**
   - Create milestone for next version
   - Prioritize features/fixes
   - Update roadmap

## Release Statistics

Track release metrics:
- Download counts per platform
- Installation method popularity
- User feedback and issues
- Performance metrics

## Security

### Signing Binaries (Future Enhancement)

Consider adding binary signing:
```yaml
# In .goreleaser.yml
signs:
  - artifacts: checksum
    args:
      - "--batch"
      - "--local-user"
      - "{{ .Env.GPG_FINGERPRINT }}"
      - "--output"
      - "${signature}"
      - "--detach-sign"
      - "${artifact}"
```

### Supply Chain Security

- All dependencies are vendored
- Builds run in isolated GitHub Actions runners
- Checksums provided for verification
- Source code is publicly auditable

## Rollback Procedure

If a release has critical issues:

1. **Mark as Pre-release**
   - Edit the release on GitHub
   - Check "This is a pre-release"

2. **Create Hotfix Release**
   ```bash
   git tag -a v1.0.1 -m "Hotfix: Critical bug fix"
   git push origin v1.0.1
   ```

3. **Notify Users**
   - Post announcement
   - Update documentation
   - Provide upgrade instructions

## Continuous Improvement

Future enhancements:
- [ ] Add binary signing (GPG/Cosign)
- [ ] Implement SBOM generation
- [ ] Add vulnerability scanning
- [ ] Create Docker images
- [ ] Add Scoop/Chocolatey support
- [ ] Implement auto-update mechanism
- [ ] Add telemetry for usage statistics

## Resources

- **GitHub Actions**: https://docs.github.com/en/actions
- **GoReleaser**: https://goreleaser.com/
- **Semantic Versioning**: https://semver.org/
- **Release Best Practices**: https://docs.github.com/en/repositories/releasing-projects-on-github/about-releases

---

**Last Updated**: April 16, 2026  
**Version**: 1.0.0  
**Maintainer**: DevPlatform Team
