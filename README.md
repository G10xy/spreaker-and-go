# spreaker-cli

A command-line interface for the [Spreaker](https://www.spreaker.com) podcast platform.

## Overview

`spreaker-cli` allows you to manage your Spreaker podcasts directly from the terminal. You can manage shows, upload episodes, view statistics, search content, and more — without opening a browser.

## Features

- **Show Management** — Create, list, update, and delete podcast shows
- **Episode Management** — Upload, update, download, and manage episodes
- **Statistics** — View plays, likes, followers, geographic data, and more
- **Search** — Find shows and episodes across the platform
- **Social Features** — Follow users, like episodes, manage favorites
- **Multiple Output Formats** — Table (human-friendly), JSON (scripting), or plain text
- **Cross-Platform** — Works on Linux, macOS, and Windows

## Installation

Download the latest release for your platform from [GitHub Releases](https://github.com/G10xy/spreaker-and-go/releases).

```bash
# Linux
chmod +x spreaker-linux-amd64
sudo mv spreaker-linux-amd64 /usr/local/bin/spreaker

# macOS (Apple Silicon)
chmod +x spreaker-macos-arm64
sudo mv spreaker-macos-arm64 /usr/local/bin/spreaker

# Verify installation
spreaker --version
```

## Quick Start

```bash
# 1. Login with your API token
spreaker login

# 2. View your profile
spreaker me

# 3. List your shows
spreaker shows list

# 4. List episodes
spreaker episodes list <show-id>
```

To get an API token, visit [Spreaker Developer Settings](https://www.spreaker.com/account/developers).

## Documentation

Full documentation is available in the [`docs/`](docs/) folder:

- [Getting Started](docs/getting-started.md) — Installation, authentication, configuration
- [Users](docs/users.md) — User profiles, followers, blocking
- [Shows](docs/shows.md) — Show management and favorites
- [Episodes](docs/episodes.md) — Episode management, likes, bookmarks
- [Messages](docs/messages.md) — Episode comments
- [Chapters](docs/chapters.md) — Episode chapters
- [Cuepoints](docs/cuepoints.md) — Ad injection points
- [Statistics](docs/statistics.md) — Analytics and metrics
- [Search](docs/search.md) — Search shows and episodes
- [Explore](docs/explore.md) — Browse by category
- [Tags](docs/tags.md) — Discover by tags
- [Miscellaneous](docs/miscellaneous.md) — Categories and languages

## Command Overview

```
spreaker
├── login                 # Authenticate with API token
├── me                    # View your profile
├── users                 # Manage users (get, follow, block, etc.)
├── shows                 # Manage shows (list, create, update, delete, favorites)
├── episodes              # Manage episodes (list, upload, update, download, likes)
├── stats                 # View statistics (plays, likes, geo, devices, etc.)
├── search                # Search shows and episodes
├── explore               # Browse shows by category
├── tags                  # Find episodes by tag
├── chapters              # Manage episode chapters
├── cuepoints             # Manage ad cuepoints
├── messages              # Manage episode messages
├── misc                  # List categories and languages
└── config                # Manage CLI configuration
```

## Output Formats

All commands support three output formats via `--output` (or `-o`):

| Format | Description | Use Case |
|--------|-------------|----------|
| `table` | Aligned columns (default) | Human reading |
| `json` | JSON array/object | Scripting, piping |
| `plain` | Tab-separated values | Shell scripts |

## Configuration

Configuration is stored in `~/.config/spreaker-cli/config.yaml` (Linux/macOS) or `%APPDATA%\spreaker-cli\config.yaml` (Windows).

```bash
# View configuration
spreaker config show

# Set default show
spreaker config set default_show_id 12345

# Set default output format
spreaker config set output_format json
```

Environment variables can override config: `SPREAKER_TOKEN`, `SPREAKER_OUTPUT`.

## Building from Source

Requires Go 1.21 or later.

```bash
# Clone and build
git clone https://github.com/G10xy/spreaker-and-go.git
cd spreaker-and-go
go build -o spreaker ./cmd/spreaker

# Run tests
go test ./...
```

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## License

MIT License — see [LICENSE](LICENSE) file for details.

## Links

- [Spreaker](https://www.spreaker.com)
- [Spreaker API Documentation](https://developers.spreaker.com/api/)
- [Report Issues](https://github.com/G10xy/spreaker-and-go/issues)
