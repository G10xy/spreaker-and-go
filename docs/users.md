# Users

Manage user profiles and social interactions.

API Reference: https://developers.spreaker.com/api/users/

## Commands

### me

Display the authenticated user's profile.

```bash
spreaker me
```

### users get

Get a user's public profile by ID.

```bash
spreaker users get <user-id>
```

### users update

Update your profile.

```bash
spreaker users update --fullname "John Doe"
spreaker users update --description "Podcast enthusiast"
spreaker users update --username johndoe
```

| Flag | Description |
|------|-------------|
| `--fullname` | Display name |
| `--description` | Bio/description |
| `--username` | Username |
| `--gender` | Gender (male, female, other) |
| `--birthday` | Birthday (YYYY-MM-DD) |
| `--location` | Location |
| `--contact-email` | Contact email |

### users shows

List a user's shows.

```bash
spreaker users shows <user-id>
spreaker users shows <user-id> --limit 50
```

### users followers

List a user's followers.

```bash
spreaker users followers <user-id>
spreaker users followers <user-id> --limit 50
```

### users followings

List who a user follows.

```bash
spreaker users followings <user-id>
spreaker users followings <user-id> --limit 50
```

### users follow

Follow a user.

```bash
spreaker users follow <user-id>
```

### users unfollow

Unfollow a user.

```bash
spreaker users unfollow <user-id>
```

### users blocks

List your blocked users.

```bash
spreaker users blocks
spreaker users blocks --limit 50
```

### users block

Block a user.

```bash
spreaker users block <user-id>
```

### users unblock

Unblock a user.

```bash
spreaker users unblock <user-id>
```
