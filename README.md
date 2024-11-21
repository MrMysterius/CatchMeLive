# CatchMeLive
> Minimal Space Usage Live Checker and Discord Notifier for Twitch written in GO

Checks the live status of a Twitch Channel and if live sends a Discord Message through a Webhook which will get also get updated periodically.

## Prerequisites

- Knowledge of Docker
- You will need a Confidential Application on Twitch https://dev.twitch.tv/console
- Discord Webhook

## Installation / Usage

- Clone this Repo `git clone https://github.com/MrMysterius/CatchMeLive.git`
- Edit the docker-compose.yml to include all the required environment variable values (TWITCH_CLIENT_ID, TWITCH_CLIENT_SECRET, TWITCH_CHANNEL_NAME, DISCORD_WEBHOOK_URL)
- Launch it `docker compose up -d`

## Live Message

The live message can include the following placeholders:

PLACEHOLDER | DESCRIPTION
--- | ---
{CHANNEL_NAME} | Replaces the text with the name of the channel
{GAME} | Replaces the text with the name of the game/category of the stream
{TITLE} | Replaces the text with the title of the stream
