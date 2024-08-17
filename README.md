
# Pigeon Box

Pigeon box is a simple, secure, open-source chat application built on top of a Slack workspace bot.

Slack is used only for chat(thread) initialization and user authentication, hence Slack never sees the messages or files shared.


---

## Motivation

1. Secure communication for my team, eliminating the need for matrix,pastebin etc.

2. Recent Slack oopses:
    ##### Slack(Salesforce) wants to use your business data for their AI/ML model training
    > To develop AI/ML models, our systems analyze Customer Data (e.g. messages, content and files) submitted to Slack.[^1]
    
    ##### Disney's Slack data leaked
    > The data allegedly includes every message and file from nearly 10,000 channels, including unreleased projects, code, images, login credentials, and links to internal websites and APIs.[^2]


[^1]: https://www.theregister.com/2024/05/20/slack_ts_and_cs_update/

[^2]: https://www.wired.com/story/disney-slack-leak-nullbulge/


---

## How it works

1. Slack user creates a thread via the bot command.
2. Bot creates a message inviting the other slack group user(s) to the thread.
3. Users are able to request a one time link to access the thread.
4. After authenticating, they're able to share messages and files securely.

Messages are stored on your server and are deleted after the set expiration time. They're encrypted using thread specific keys and are never visible to Slack.

<details>
  <summary>User flow visualized</summary>

![flow](docs/pigeonbox-flow.jpg "Pigeon Box Flow")

</details>

<details>
  <summary>Encryption visualized</summary>

![encryption](docs/pigeonbox-encryption.jpg "Pigeon Box Encryption")

</details>

---

## Deployment Guides

1. [Slack Bot](docs/GUIDES.md#slack-bot)
2. [Environment Variables](docs/GUIDES.md#environment-variables)
3. [Database](docs/GUIDES.md#database)
4. [Server](docs/GUIDES.md#server)
5. [File Storage](docs/GUIDES.md#file-storage)


---


## Goals

1. Keeping a good balance between security and usability.
    - Should be accessible to non-technical users.
    - Should be secure enough to be used by security-conscious teams.
2. Making the app very cost-effective to run.
    - Should be able to run on a single shared instance.
    - Should be able to run for free or a sub $2/month budget.
    - Should be able to run on your existing infrastructure.
3. Crafting a beautiful, responsive and accessible UI.
    - Should be very easy to use on desktop and mobile.
    - Should be very intuitive to use and pleasing to the eye.
    - Should be accessible to screen readers and keyboard users.
    - Should be very light on user resources.


---




