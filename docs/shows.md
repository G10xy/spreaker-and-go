# Shows

Manage your podcast shows.

API Reference: https://developers.spreaker.com/api/shows/

## Commands

### shows list

List all your shows.

```bash
spreaker shows list
spreaker shows list --limit 50
```

### shows get

Get details of a specific show.

```bash
spreaker shows get <show-id>
```

### shows create

Create a new show.

```bash
spreaker shows create --title "My Podcast"
spreaker shows create --title "My Podcast" --language en --category 1
```

| Flag | Description |
|------|-------------|
| `--title` | Show title (required) |
| `--description` | Show description |
| `--language` | Language code (e.g., en, it, es) |
| `--category` | Category ID |
| `--explicit` | Mark as explicit content |

### shows update

Update an existing show.

```bash
spreaker shows update <show-id> --title "New Title"
spreaker shows update <show-id> --description "New description"
```

| Flag | Description |
|------|-------------|
| `--title` | Show title |
| `--description` | Show description |
| `--language` | Language code (e.g., en, it, es) |
| `--category` | Category ID |
| `--explicit` | Mark as explicit content |

### shows delete

Delete a show permanently.

```bash
spreaker shows delete <show-id>
spreaker shows delete <show-id> --force
```

| Flag | Description |
|------|-------------|
| `--force`, `-f` | Skip confirmation prompt |

### shows favorites

List your favorite shows.

```bash
spreaker shows favorites
spreaker shows favorites --limit 50
```

### shows favorite

Add a show to your favorites.

```bash
spreaker shows favorite <show-id>
```

### shows unfavorite

Remove a show from your favorites.

```bash
spreaker shows unfavorite <show-id>
```
