# Guides

Table of Contents

- [Environment Variables](#environment-variables)
- [Turso](#turso)
- [Fly.io](#flyio)
  - [Steps to deploy](#steps-to-deploy)
  - [Setting up secrets](#setting-up-secrets)
  
  

## Slack Bot


<details>
  <summary>App Manifest</summary>

``` yaml
{
    "display_information": {
        "name": "Pigeon box",
        "description": "Communicate securely within your organization",
        "background_color": "#050405"
    },
    "features": {
        "app_home": {
            "home_tab_enabled": false,
            "messages_tab_enabled": true,
            "messages_tab_read_only_enabled": false
        },
        "bot_user": {
            "display_name": "Pigeon box",
            "always_online": true
        },
        "slash_commands": [
            {
                "command": "/pigeon",
                "description": "Create new thread",
                "usage_hint": " [optional name]",
                "should_escape": false
            }
        ]
    },
    "oauth_config": {
        "scopes": {
            "bot": [
                "app_mentions:read",
                "channels:history",
                "channels:join",
                "channels:read",
                "chat:write",
                "chat:write.customize",
                "commands",
                "groups:history",
                "groups:read",
                "im:history",
                "im:read",
                "im:write",
                "mpim:read",
                "mpim:write",
                "reactions:read",
                "reactions:write",
                "team:read",
                "usergroups:read",
                "users.profile:read",
                "users:read"
            ]
        }
    },
    "settings": {
        "event_subscriptions": {
            "bot_events": [
                "app_mention",
                "member_joined_channel",
                "message.channels",
                "message.im"
            ]
        },
        "interactivity": {
            "is_enabled": true
        },
        "org_deploy_enabled": false,
        "socket_mode_enabled": true,
        "token_rotation_enabled": false
    }
}
```
</details>



## Environment Variables

Create a `.env` from the `.env.example` file.
```bash
cp .env.example .env
```

## Turso



## Fly.io

#### Steps to deploy

1. Install the fly CLI (https://fly.io/docs/flyctl/install/) and login.
2. Launch the app from fly.toml, by running `fly launch`
3. Set the env variables in your local .env file (see .env.example)
4. Set secrets using the script `fly_secrets.sh` (see below)


#### Setting up secrets

Assuming you have the fly CLI installed and logged in, and you have env variables set in `.env` file, you can use the following script to set secrets for your fly.io app.

```bash
chmod +x fly_secrets.sh
./fly_secrets.sh
```

