# Explore

Discover podcasts by browsing categories.

API Reference: https://developers.spreaker.com/api/explore/

## Commands

### explore category

List shows in a specific category, ranked by popularity and quality.

```bash
spreaker explore category <category-id>
spreaker explore category 14 --limit 50
```

| Flag | Description |
|------|-------------|
| `--limit`, `-l` | Maximum number of shows (default: 20) |

Use `spreaker misc categories` to see available category IDs.
