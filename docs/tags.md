# Tags

Discover episodes by searching for specific tags/hashtags.

API Reference: https://developers.spreaker.com/api/tags/

## Commands

### tags episodes

Get the latest episodes with a specific tag.

```bash
spreaker tags episodes <tag-name>
spreaker tags episodes "breaking news"
spreaker tags episodes tech
spreaker tags episodes "machine learning" --limit 50
```

The tag name can contain spaces and special characters.

| Flag | Description |
|------|-------------|
| `--limit`, `-l` | Maximum number of episodes (default: 20) |
