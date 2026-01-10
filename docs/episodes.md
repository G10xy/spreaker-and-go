# Episodes

Manage podcast episodes.

API Reference: https://developers.spreaker.com/api/episodes/

## Commands

### episodes list

List episodes of a show.

```bash
spreaker episodes list                    # Uses default show from config
spreaker episodes list <show-id>
spreaker episodes list <show-id> --limit 50
```

### episodes get

Get details of a specific episode.

```bash
spreaker episodes get <episode-id>
```

### episodes upload

Upload a new episode to a show.

```bash
spreaker episodes upload <show-id> <audio-file> --title "Episode Title"
spreaker episodes upload <show-id> ./episode.mp3 --title "Episode 1" --description "First episode" --tags "news,tech"
```

| Flag | Description |
|------|-------------|
| `--title`, `-t` | Episode title (required) |
| `--description`, `-d` | Episode description |
| `--tags` | Tags (comma-separated) |
| `--explicit` | Mark as explicit content |
| `--downloadable` | Allow downloads (default: true) |

### episodes update

Update an existing episode.

```bash
spreaker episodes update <episode-id> --title "New Title"
spreaker episodes update <episode-id> --description "Updated description"
spreaker episodes update <episode-id> --hidden
```

| Flag | Description |
|------|-------------|
| `--title` | Episode title |
| `--description` | Episode description |
| `--tags` | Tags (comma-separated) |
| `--explicit` | Mark as explicit content |
| `--downloadable` | Allow downloads |
| `--hidden` | Hide the episode |

### episodes draft

Create a draft episode without an audio file.

```bash
spreaker episodes draft <show-id> --title "Upcoming Episode"
spreaker episodes draft <show-id> --title "Draft" --description "Work in progress"
```

| Flag | Description |
|------|-------------|
| `--title` | Episode title (required) |
| `--description` | Episode description |
| `--tags` | Tags (comma-separated) |
| `--explicit` | Mark as explicit content |
| `--downloadable` | Allow downloads (default: true) |
| `--hidden` | Hide the episode |

### episodes delete

Delete an episode permanently.

```bash
spreaker episodes delete <episode-id>
spreaker episodes delete <episode-id> --force
```

| Flag | Description |
|------|-------------|
| `--force`, `-f` | Skip confirmation prompt |

### episodes download

Download an episode's audio file.

```bash
spreaker episodes download <episode-id>
spreaker episodes download <episode-id> --output ~/podcasts/episode.mp3
spreaker episodes download <episode-id> --url-only
```

| Flag | Description |
|------|-------------|
| `--output`, `-O` | Output file path (default: episode title) |
| `--url-only`, `-u` | Only print the download URL |

### episodes download-all

Download all episodes of a show. Files that already exist are skipped by default (resume capability).

```bash
spreaker episodes download-all <show-id>
spreaker episodes download-all <show-id> --output-dir ~/podcasts/myshow
spreaker episodes download-all <show-id> --limit 10
spreaker episodes download-all <show-id> --no-skip-existing
```

| Flag | Description |
|------|-------------|
| `--output-dir`, `-O` | Output directory (default: ./<show-title>/) |
| `--skip-existing` | Skip episodes that already exist (default: true) |
| `--limit`, `-l` | Maximum number of episodes to download (0 = all) |

### episodes likes

List your liked episodes.

```bash
spreaker episodes likes
spreaker episodes likes --limit 50
```

### episodes like

Like an episode.

```bash
spreaker episodes like <episode-id>
```

### episodes unlike

Unlike an episode.

```bash
spreaker episodes unlike <episode-id>
```

### episodes bookmark

Bookmark an episode.

```bash
spreaker episodes bookmark <episode-id>
```

### episodes unbookmark

Remove an episode from bookmarks.

```bash
spreaker episodes unbookmark <episode-id>
```
