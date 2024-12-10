![Voice Watch](/assets/banner.png)
![GitHub Release](https://img.shields.io/github/v/release/jord-nijhuis/discord-voice-watch)
![GitHub License](https://img.shields.io/github/license/jord-nijhuis/discord-voice-watch)

> A Discord bot that notifies users when someone joins a voice channel

## Usage

To add the bot to your server, follow this [link](https://discord.com/oauth2/authorize?client_id=1314977679434448968).
Users can then call `/voice-watch enable` to receive a DM when someone joins a voice channel of the server. If users
no longer would like to receive notifications, they can use `/voice-watch disable`.

## Hosting the bot yourself

To host the bot yourself, grab the latest binary from the [releases](https://github.com/jord-nijhuis/discord-voice-watch/releases/) 
page. Running the bot for the first time will create a `config.yml` file in the working directory. This file should look
like this:

```yaml
discord:
  token: "[TOKEN]" # The token of the discord bot
logging:
  level: warning # The log level of the bot
notifications:
    delay-before-sending: 1m # How long a user should be in a voice channel before a notification is sent
    delay-between-messages: 1h # How long the bot should wait before sending another notification to the same user
```

Be sure to replace `[TOKEN]` with the token of your bot. You can create a bot in the Discord Developer Portal. The
bot requires the following permissions to function:

- View Channels
- Send Messages

