# spreaker-cli

A command-line interface for the [Spreaker](https://www.spreaker.com) podcast platform.

## Overview

`spreaker-cli` allows you to manage your Spreaker podcasts directly from the terminal. You can list shows, upload episodes, view statistics, and more — without opening a browser.

## Features

- **Authentication** — Login once, token is saved securely
- **Show Management** — List, view, and delete your podcast shows
- **Episode Management** — List, view, upload, and delete episodes
- **Multiple Output Formats** — Table (human-friendly), JSON (scripting), or plain text
- **Cross-Platform** — Works on Linux, macOS, and Windows

## Installation

### From Source

Requires Go 1.21 or later.

```bash
git clone https://github.com/G10xy/spreaker-and-go.git
cd spreaker-cli
go build -o spreaker ./cmd/spreaker

# Optional: move to PATH
sudo mv spreaker /usr/local/bin/
```

### Verify Installation

```bash
spreaker --version
```

## Getting Started

### 1. Get Your API Token

1. Go to [Spreaker Developer Settings](https://www.spreaker.com/account/developers)
2. Register a new application (use `http://localhost:8080/callback` as callback URL)
3. Generate an access token

### 2. Login

```bash
spreaker login
# Enter your token when prompted
```

Your token is saved to `~/.config/spreaker-and-go/config.yaml`.

### 3. Start Using

```bash
# View your profile
spreaker me

# List your shows
spreaker shows list

# List episodes of a show
spreaker episodes list 12345
```

## Usage

### Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--output` | `-o` | Output format: `table`, `json`, `plain` |
| `--token` | | Override saved token for this command |
| `--help` | `-h` | Show help |
| `--version` | `-v` | Show version |

### Commands

#### Authentication

```bash
# Login and save token
spreaker login

# View current user
spreaker me
```

#### Shows

```bash
# List all your shows
spreaker shows list

# List with limit
spreaker shows list --limit 50

# Get show details
spreaker shows get <show-id>

# Delete a show (with confirmation)
spreaker shows delete <show-id>

# Delete without confirmation
spreaker shows delete <show-id> --force
```

#### Episodes

```bash
# List episodes (uses default show if configured)
spreaker episodes list

# List episodes of a specific show
spreaker episodes list <show-id>

# Get episode details
spreaker episodes get <episode-id>

# Upload a new episode
spreaker episodes upload <show-id> ./episode.mp3 --title "Episode Title"

# Upload with full options
spreaker episodes upload <show-id> ./episode.mp3 \
  --title "Episode 42: The Answer" \
  --description "In this episode we discuss everything." \
  --tags "science,philosophy" \
  --explicit

# Delete an episode
spreaker episodes delete <episode-id>
```

#### Configuration

```bash
# Show current configuration
spreaker config show

# Show config file path
spreaker config path

# Set default show (so you don't have to specify it every time)
spreaker config set default_show_id 12345

# Set default output format
spreaker config set output_format json
```

### Output Formats

**Table (default)** — Human-readable aligned columns:

```
$ spreaker shows list
ID       TITLE                EPISODES  FOLLOWERS  PLAYS
--       -----                --------  ---------  -----
123456   My Tech Podcast      45        1200       50000
789012   Weekly News          120       3500       125000
```

**JSON** — For scripting and piping to other tools:

```bash
$ spreaker shows list --output json
[
  {
    "show_id": 123456,
    "title": "My Tech Podcast",
    ...
  }
]
```

**Plain** — Minimal output, one record per line:

```bash
$ spreaker shows list --output plain
123456	My Tech Podcast
789012	Weekly News
```

### Examples

```bash
# Get your latest show's episodes as JSON
spreaker episodes list $(spreaker shows list -o plain | head -1 | cut -f1) -o json

# Upload and get the new episode ID
spreaker episodes upload 12345 ./ep1.mp3 --title "Pilot" -o plain | cut -f1

# Quick check if logged in
spreaker me -o plain && echo "Logged in!" || echo "Not logged in"
```

## Configuration

Configuration is stored in:

| OS | Path |
|----|------|
| Linux | `~/.config/spreaker-cli/config.yaml` |
| macOS | `~/Library/Application Support/spreaker-cli/config.yaml` |
| Windows | `%APPDATA%\spreaker-cli\config.yaml` |

### Config Options

| Key | Description | Default |
|-----|-------------|---------|
| `token` | Your OAuth access token | (none) |
| `default_show_id` | Default show for episode commands | (none) |
| `output_format` | Default output format | `table` |
| `api_url` | API base URL | `https://api.spreaker.com` |

### Environment Variables

You can override config with environment variables:

```bash
SPREAKER_TOKEN=xxx spreaker me
SPREAKER_OUTPUT=json spreaker shows list
```

### Building

```bash
# Development build
go build -o spreaker ./cmd/spreaker

# Release build with version
go build -ldflags "-X main.version=1.0.0" -o spreaker ./cmd/spreaker

# Run tests
go test ./...
```

### Dependencies

- [cobra](https://github.com/spf13/cobra) — CLI framework
- [viper](https://github.com/spf13/viper) — Configuration management

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License — see [LICENSE](LICENSE) file for details.

## Links

- [Spreaker](https://www.spreaker.com)
- [Spreaker API Documentation](https://developers.spreaker.com/api/)
- [Report Issues](https://github.com/G10xy/spreaker-and-go/issues)
