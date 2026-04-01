# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

```bash
# Development build
go build -o spreaker ./cmd/spreaker

# Release build with version
go build -ldflags "-X main.version=1.0.0" -o spreaker ./cmd/spreaker

# Run tests
go test ./...

# Cross-compile (examples)
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=v1.0.0" -o spreaker-linux-amd64 ./cmd/spreaker
GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=v1.0.0" -o spreaker-macos-arm64 ./cmd/spreaker
```

## Architecture

This is a CLI application for the Spreaker podcast platform API, built with Cobra and Viper.

### Package Structure

- `cmd/spreaker/main.go` - Entry point, passes version to CLI
- `internal/cli/` - Cobra command definitions (one file per command group)
- `internal/api/` - Spreaker API client with typed methods for each endpoint
- `internal/config/` - Viper-based configuration (token, defaults, output format)
- `internal/output/` - Multi-format output (table/json/plain) via `Formatter`
- `pkg/models/` - API response types with JSON tags

### Key Patterns

**Adding a new command:**
1. Create `internal/cli/<resource>.go` with `newResourceCmd()` function
2. Register in `internal/cli/root.go` via `cmd.AddCommand()`
3. Use `getClient(cmd)` and `getFormatter(cmd)` helpers from `helpers.go`

**Adding a new API endpoint:**
1. Add method to `internal/api/<resource>.go` using `Client.Get/Post/Delete` helpers
2. Create response types in `pkg/models/<resource>.go`
3. Add print methods to `internal/output/formatter.go` for table/json/plain formats

**API response wrapper:** All Spreaker API responses are wrapped in `{"response": ...}`. The client handles this automatically via `apiResponse` struct.

**Pagination:** Use `api.GetPaginated[T]()` generic function for paginated endpoints, returns `PaginatedResult[T]` with Items and HasMore.

### Configuration Priority

Token resolution order: `--token` flag > `SPREAKER_TOKEN` env > config file (`~/.config/spreaker-cli/config.yaml`)

## API Reference

Always refer to the official Spreaker API documentation: https://developers.spreaker.com/api/
