
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

## Deployment Guides

1. [Slack Bot](docs/GUIDES.md#slack-bot)
2. [Environment Variables](docs/GUIDES.md#environment-variables)
3. [Database](docs/GUIDES.md#database)
4. [Server](docs/GUIDES.md#server)
5. [File Storage](docs/GUIDES.md#file-storage)




