
# Pigeon Box

Pigeon box is a simple, secure, open-source chat application built on top of a Slack workspace bot.

Slack is used only for chat(thread) initialization and user authentication, hence Slack never sees the messages or files shared.

---

## Motivation

1. Secure communication for my team, eliminating the need for matrix/pastebin etc.

2. Recent Slack oopses:
    ##### Slack(Salesforce) wants to use your business data for their AI/ML model training
    > To develop AI/ML models, our systems analyze Customer Data (e.g. messages, content and files) submitted to Slack.
    
    After the backlash they've tweaked their T&Cs[^1], to be used only for emoji training... Right.
    
    ##### Disney's Slack data leaked
    > The data allegedly includes every message and file from nearly 10,000 channels, including unreleased projects, code, images, login credentials, and links to internal websites and APIs.[^2]
    
    Although not to the fault of their own, the data was allegedly leaked by Disney's employee, it's still a major risk.

3. Wanted to play around with go, htmx and Slack's _new_ socket mode.

[^1]: https://www.theregister.com/2024/05/20/slack_ts_and_cs_update/

[^2]: https://www.wired.com/story/disney-slack-leak-nullbulge/


---

## How it works

1. Slack user creates a thread via the bot command or an interaction. They're encouraged to set expiration time for the thread and the messages.
2. Bot creates a message inviting the other user(s) to the thread.
3. They're able to request a one time token to access the thread.
4. Once they click the link, they're authenticated to that thread only and can chat and share files with the other team members in that Slack group.

Messages are stored on your server and are deleted after the expiration time. They're encrypted using thread specific keys and are never visible to Slack.

<details>
  <summary>User flow visualized</summary>

![flow](docs/pigeonbox-flow.jpg "Pigeon Box Flow")

</details>

<details>
  <summary>Encryption visualized</summary>

![encryption](docs/pigeonbox-encryption.jpg "Pigeon Box Encryption")

</details>


---

## Deployment



  * [Slack Bot](#slack-bot)
  * [Database](#database)
  * [File Storage](#file-storage)
  * [Server](#server)
     * [General](#general)
     * [Fly.io](#flyio)


### Slack Bot




### Database

Pigeon box uses SQLite (libsql) and is compatible with managed SQLite services like [Turso](https://turso.tech/).

I'd _recommend_ going with Turso, it's easy to set up and offers a very generous free tier (as is, it would be entirely free).

If you'd still prefer local SQLite, consider persistent storage[^3].


### File Storage

Currently, Pigeon box supports the following file storage options:

- Local (should have persistent storage)[^3]
- AWS S3 (and S3 compatible services like [Tigris](https://www.tigrisdata.com/))

I'd recommend using S3 or a similar service for your deployment
By default, all files have expiration(auto delete) period, so the cost should be insignificant.


[^3]: Pigeon box is intended for short-lived chats, so even if you don't go with persistent storage, you should be fine in most cases. Just keep in mind that all user data (db, files) will be lost on each machine restart.


### Server

Pigeon can be deployed almost anywhere.
Keeping in mind, the instance has to be publicly accessible (most of the user interactions are via browser).

#### General

Dockerfile is available, and I'll be adding more deployment instructions soon/as necessary.


#### Fly.io

This is the easiest[^4] and very cost-effective[^5] way to deploy Pigeon box.
You can run it on a single `shared-cpu-1x@256MB`[^6] instance with a minimal setup, assuming you're using managed sqlite and cloud storage for files.

[Deploying to fly.io](docs/GUIDES.md#flyio)

[^4]: Minimal technical knowledge is required to deploy on fly.io.

[^5]: Assuming smallest instance in Ashburn, Virginia (US) region, you'd be looking at [~$1.94/mo](https://fly.io/docs/about/pricing/#started-fly-machines).

[^6]: This is the smallest machine size currently available on fly.io


---

### Development Notes

Building styles (tailwind)
``` bash
npx tailwindcss -i assets/static/main.css -o assets/static/build.css --watch
```

