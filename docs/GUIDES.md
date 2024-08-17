# Guides

Table of Contents

- [Environment Variables](#environment-variables)
- [Turso](#turso)
- [Fly.io](#flyio)
  - [Steps to deploy](#steps-to-deploy)
  - [Setting up secrets](#setting-up-secrets)
  



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

