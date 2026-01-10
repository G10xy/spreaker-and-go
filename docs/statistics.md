# Statistics

View statistics for users, shows, and episodes.

API Reference: https://developers.spreaker.com/api/statistics/

All time-series commands require `--from` and `--to` flags in YYYY-MM-DD format.

## Overall Statistics

### stats me

Show your overall statistics (all-time totals).

```bash
spreaker stats me
```

### stats show

Show overall statistics for a specific show.

```bash
spreaker stats show <show-id>
```

### stats episode

Show overall statistics for a specific episode.

```bash
spreaker stats episode <episode-id>
```

## Play Statistics

### stats plays

Show play statistics for a show over time.

```bash
spreaker stats plays <show-id> --from 2024-01-01 --to 2024-01-31
spreaker stats plays <show-id> --from 2024-01-01 --to 2024-01-31 --group week
```

### stats plays-user

Show play statistics for authenticated user over time.

```bash
spreaker stats plays-user --from 2024-01-01 --to 2024-01-31
```

### stats plays-episode

Show play statistics for an episode over time.

```bash
spreaker stats plays-episode <episode-id> --from 2024-01-01 --to 2024-01-31
```

### stats shows-totals

Show play totals for each of your shows.

```bash
spreaker stats shows-totals --from 2024-01-01 --to 2024-01-31
spreaker stats shows-totals --from 2024-01-01 --to 2024-01-31 --limit 50
```

### stats episodes-totals

Show play totals for each episode in a show.

```bash
spreaker stats episodes-totals <show-id> --from 2024-01-01 --to 2024-01-31
```

## Likes Statistics

### stats likes

Show likes statistics for a show over time.

```bash
spreaker stats likes <show-id> --from 2024-01-01 --to 2024-01-31
```

### stats likes-user

Show likes statistics for authenticated user over time.

```bash
spreaker stats likes-user --from 2024-01-01 --to 2024-01-31
```

### stats likes-episode

Show likes statistics for an episode over time.

```bash
spreaker stats likes-episode <episode-id> --from 2024-01-01 --to 2024-01-31
```

## Followers Statistics

### stats followers

Show followers statistics for authenticated user over time.

```bash
spreaker stats followers --from 2024-01-01 --to 2024-01-31
```

## Sources Statistics

### stats sources

Show play/download sources for a show.

```bash
spreaker stats sources <show-id> --from 2024-01-01 --to 2024-01-31
```

### stats sources-user

Show play/download sources for authenticated user.

```bash
spreaker stats sources-user --from 2024-01-01 --to 2024-01-31
```

### stats sources-episode

Show play/download sources for an episode.

```bash
spreaker stats sources-episode <episode-id> --from 2024-01-01 --to 2024-01-31
```

## Devices Statistics

### stats devices

Show device breakdown for a show.

```bash
spreaker stats devices <show-id> --from 2024-01-01 --to 2024-01-31
```

### stats devices-user

Show device breakdown for authenticated user.

```bash
spreaker stats devices-user --from 2024-01-01 --to 2024-01-31
```

### stats devices-episode

Show device breakdown for an episode.

```bash
spreaker stats devices-episode <episode-id> --from 2024-01-01 --to 2024-01-31
```

## Operating System Statistics

### stats os

Show operating system breakdown for a show.

```bash
spreaker stats os <show-id> --from 2024-01-01 --to 2024-01-31
```

### stats os-user

Show operating system breakdown for authenticated user.

```bash
spreaker stats os-user --from 2024-01-01 --to 2024-01-31
```

### stats os-episode

Show operating system breakdown for an episode.

```bash
spreaker stats os-episode <episode-id> --from 2024-01-01 --to 2024-01-31
```

## Geographic Statistics

### stats geo

Show geographic breakdown for a show.

```bash
spreaker stats geo <show-id> --from 2024-01-01 --to 2024-01-31
```

### stats geo-user

Show geographic breakdown for authenticated user.

```bash
spreaker stats geo-user --from 2024-01-01 --to 2024-01-31
```

## Listeners Statistics

### stats listeners

Show unique listeners for a show over time.

```bash
spreaker stats listeners <show-id> --from 2024-01-01 --to 2024-01-31
```

## Common Flags

| Flag | Description |
|------|-------------|
| `--from` | Start date (YYYY-MM-DD, required) |
| `--to` | End date (YYYY-MM-DD, required) |
| `--group` | Group by: day, week, or month (default: day) |
| `--limit`, `-l` | Maximum results for totals commands |
