# Getting Started

This guide will help you install and configure spreaker-cli to manage your Spreaker podcasts from the terminal.

## Installation

Download the latest release for your platform from [GitHub Releases](https://github.com/G10xy/spreaker-and-go/releases).

**Linux:**
```bash
chmod +x spreaker-linux-amd64
sudo mv spreaker-linux-amd64 /usr/local/bin/spreaker
```

**macOS (Apple Silicon):**
```bash
chmod +x spreaker-macos-arm64
sudo mv spreaker-macos-arm64 /usr/local/bin/spreaker
```

**macOS (Intel):**
```bash
chmod +x spreaker-macos-amd64
sudo mv spreaker-macos-amd64 /usr/local/bin/spreaker
```

**Windows:**
- Download `spreaker-windows-amd64.exe`
- Rename to `spreaker.exe`
- Add to your PATH

### Verify Installation

```bash
spreaker --version
```

## Authentication

### Get Your API Token

1. Go to [Spreaker Developer Settings](https://www.spreaker.com/account/developers)
2. Register a new application (use `http://localhost:8080/callback` as callback URL)
3. Generate an access token

### Login

```bash
spreaker login
```

You will be prompted to enter your API token. The CLI validates the token by calling the API, then saves it to your config file.

```
Enter your Spreaker API token: <paste-your-token>
Logged in as John Doe (@johndoe)
  Token saved to /home/user/.config/spreaker-cli/config.yaml
```

### Verify Authentication

```bash
spreaker me
```

This displays your profile information and confirms authentication is working.

## Configuration

Configuration is stored in:

| OS | Path |
|----|------|
| Linux | `~/.config/spreaker-cli/config.yaml` |
| macOS | `~/Library/Application Support/spreaker-cli/config.yaml` |
| Windows | `%APPDATA%\spreaker-cli\config.yaml` |

### View Current Configuration

```bash
spreaker config show
```

### Set Default Show

To avoid specifying a show ID for every episode command:

```bash
spreaker config set default_show_id 12345
```

### Set Default Output Format

```bash
spreaker config set output_format json
```

Available formats: `table` (default), `json`, `plain`

### Environment Variables

Override configuration with environment variables:

```bash
SPREAKER_TOKEN=xxx spreaker me
SPREAKER_OUTPUT=json spreaker shows list
```

## Output Formats

All commands support three output formats via the `--output` (or `-o`) flag:

### Table (default)

Human-readable aligned columns:

```
$ spreaker shows list
ID       TITLE                EPISODES  FOLLOWERS  PLAYS
--       -----                --------  ---------  -----
123456   My Tech Podcast      45        1200       50000
```

### JSON

Machine-readable, ideal for scripting:

```bash
$ spreaker shows list --output json
```

### Plain

Tab-separated, one record per line:

```bash
$ spreaker shows list --output plain
123456	My Tech Podcast
```

## Global Flags

These flags are available on all commands:

| Flag | Short | Description |
|------|-------|-------------|
| `--output` | `-o` | Output format: `table`, `json`, `plain` |
| `--token` | | Override saved token for this command |
| `--help` | `-h` | Show help |
| `--version` | `-v` | Show version |
