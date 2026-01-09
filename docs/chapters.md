# Chapters

Manage episode chapters. Chapters are bookmarks within an episode that help listeners navigate to specific points.

API Reference: https://developers.spreaker.com/api/episodes/chapters/

## Commands

### chapters list

List all chapters for an episode.

```bash
spreaker chapters list <episode-id>
spreaker chapters list <episode-id> --limit 50
```

Aliases: `chapter list`

### chapters add

Add a new chapter to an episode.

```bash
spreaker chapters add <episode-id> --starts-at 30000 --title "Introduction"
spreaker chapters add <episode-id> --starts-at 120000 --title "Main Topic" --url "https://example.com"
```

| Flag | Description |
|------|-------------|
| `--starts-at` | Position in milliseconds (required) |
| `--title` | Chapter title, max 120 chars (required) |
| `--url` | External URL for extra information |
| `--image` | Image file path (400x400+, max 5MB, JPG/PNG) |
| `--crop` | Crop coordinates: x1,y1,x2,y2 |

### chapters update

Update an existing chapter.

```bash
spreaker chapters update <episode-id> <chapter-id> --title "New Title"
spreaker chapters update <episode-id> <chapter-id> --starts-at 60000
spreaker chapters update <episode-id> <chapter-id> --image remove
```

| Flag | Description |
|------|-------------|
| `--starts-at` | Position in milliseconds |
| `--title` | Chapter title |
| `--url` | External URL |
| `--image` | Image file path (or 'remove' to delete) |
| `--crop` | Crop coordinates: x1,y1,x2,y2 |

### chapters delete

Delete a single chapter from an episode.

```bash
spreaker chapters delete <episode-id> <chapter-id>
```

### chapters delete-all

Delete all chapters from an episode.

```bash
spreaker chapters delete-all <episode-id>
```
