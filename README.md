# Laravel Forge Deployment CLI

An interactive command-line tool written in Go that generates deployment configurations and workflows for the [deploy-to-laravel-forge](https://github.com/the-trybe/deploy-to-laravel-forge) action.

## Installation

### Download Pre-built Binary

Download the latest release for your platform:

**Linux (amd64):**

```bash
curl -L https://github.com/the-trybe/forge-deploy-cli/releases/latest/download/forge-deploy-linux-amd64 -o forge-deploy
chmod +x forge-deploy
sudo mv forge-deploy /usr/local/bin/
```

**macOS (Intel):**

```bash
curl -L https://github.com/the-trybe/forge-deploy-cli/releases/latest/download/forge-deploy-darwin-amd64 -o forge-deploy
chmod +x forge-deploy
sudo mv forge-deploy /usr/local/bin/
```

**macOS (Apple Silicon):**

```bash
curl -L https://github.com/the-trybe/forge-deploy-cli/releases/latest/download/forge-deploy-darwin-arm64 -o forge-deploy
chmod +x forge-deploy
sudo mv forge-deploy /usr/local/bin/
```

**Windows:**
I don't care, build from source (if it works).

### Build from Source

Requires Go 1.23+:

```bash
git clone https://github.com/the-trybe/forge-deploy-cli.git
cd deploy-to-forge-cli
make build
```

## Quick Start

```bash
forge-deploy generate
```

## Usage

```bash
forge-deploy generate [options]
```

Options:

- `-f`, `--forge-config` string Forge deployment config filename (default "forge-deploy.yml")
- `-h`, `--help` help for generate
- `-o`, `--output-dir` string Output directory for generated files (default ".")
- `-b`, `--trigger-branch` string Branch that triggers deployment (default "main")
- `-w`, `--workflow-file` string GitHub Actions workflow filename (default "deploy.yml")

## Generated Files

The tool generates 4 files:

1. **forge-deploy.yml** - Declarative Forge configuration
2. **.github/workflows/deploy.yml** - GitHub Actions workflow

## Requirements

- **Runtime:** None (compiled binary)
- **Build:** Go 1.23+ (only for building from source)

## References

- [GitHub Action Documentation](https://github.com/the-trybe/deploy-to-laravel-forge)

## License

MIT License
