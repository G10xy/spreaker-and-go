# Search

Search for shows and episodes on Spreaker.

API Reference: https://developers.spreaker.com/api/search/

## Commands

### search shows

Search for shows globally.

```bash
spreaker search shows "tech podcast"
spreaker search shows "comedy" --limit 50
spreaker search shows "news" --filter editable
```

| Flag | Description |
|------|-------------|
| `--limit`, `-l` | Maximum number of results (default: 20) |
| `--filter` | Filter: `listenable` (default) or `editable` |

### search episodes

Search for episodes globally.

```bash
spreaker search episodes "artificial intelligence"
spreaker search episodes "interview" --limit 50
```

| Flag | Description |
|------|-------------|
| `--limit`, `-l` | Maximum number of results (default: 20) |
| `--filter` | Filter: `listenable` (default) or `editable` |

### search user-shows

Search for shows by a specific user.

```bash
spreaker search user-shows <user-id> "podcast"
spreaker search user-shows 12345 "tech" --limit 50
```

| Flag | Description |
|------|-------------|
| `--limit`, `-l` | Maximum number of results (default: 20) |
| `--filter` | Filter: `listenable` (default) or `editable` |

### search user-episodes

Search for episodes by a specific user.

```bash
spreaker search user-episodes <user-id> "bonus"
spreaker search user-episodes 12345 "interview" --limit 50
```

| Flag | Description |
|------|-------------|
| `--limit`, `-l` | Maximum number of results (default: 20) |
| `--filter` | Filter: `listenable` (default) or `editable` |

### search show-episodes

Search for episodes within a specific show.

```bash
spreaker search show-episodes <show-id> "bonus episode"
spreaker search show-episodes 12345 "guest" --limit 50
```

| Flag | Description |
|------|-------------|
| `--limit`, `-l` | Maximum number of results (default: 20) |
| `--filter` | Filter: `listenable` (default) or `editable` |
