# Messages

Manage episode messages. Messages are public comments that listeners can leave on episodes.

API Reference: https://developers.spreaker.com/api/episodes/messages/

## Commands

### messages list

List all messages for an episode.

```bash
spreaker messages list <episode-id>
spreaker messages list <episode-id> --limit 50
```

Aliases: `msg list`

### messages create

Leave a message on an episode. Maximum 4000 characters.

```bash
spreaker messages create <episode-id> "Great episode!"
spreaker messages create <episode-id> "Thanks for the insights!"
```

Aliases: `msg create`

### messages delete

Delete a message from an episode. You can only delete your own messages or messages on your episodes.

```bash
spreaker messages delete <episode-id> <message-id>
```

Aliases: `msg delete`

### messages report

Report a message as spam or abuse. Reported messages are reviewed by Spreaker staff.

```bash
spreaker messages report <episode-id> <message-id>
```

Aliases: `msg report`
