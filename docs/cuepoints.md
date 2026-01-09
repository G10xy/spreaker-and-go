# Cuepoints

Manage episode cuepoints for ad injection. Cuepoints are specific points in time within an episode where audio ads can be injected.

API Reference: https://developers.spreaker.com/api/episodes/cuepoints/

Note: Setting cuepoints alone is not enough to get ads injected. You also need to enable Ads and Monetization capabilities on your account and show.

## Commands

### cuepoints list

List all cuepoints for an episode.

```bash
spreaker cuepoints list <episode-id>
```

Aliases: `cue list`

### cuepoints set

Set cuepoints for an episode. This replaces all existing cuepoints.

Format: `timecode:max_ads` where timecode is in milliseconds.

```bash
# Set a single cuepoint at 30 seconds with max 1 ad
spreaker cuepoints set <episode-id> 30000:1

# Set multiple cuepoints
spreaker cuepoints set <episode-id> 30000:1 60000:2 90000:1

# Clear all cuepoints (set empty list)
spreaker cuepoints set <episode-id>
```

Aliases: `cue set`

### cuepoints delete

Delete all cuepoints for an episode.

```bash
spreaker cuepoints delete <episode-id>
```

Aliases: `cue delete`
