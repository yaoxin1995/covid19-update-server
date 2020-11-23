# Notification Service

This web service uses Telegram bot API to send messages.
TODO...

## Prerequisites

In order to be able to send messages via Telegram you need create a
[Telegram bot](https://telegram.org/blog/bot-revolution).
Read [this](https://core.telegram.org/bots) introduction to get in touch with the Telegram bot platform.
After you created a bot you get a bot token.


## How to run the web service

 * open `docker-compose.yml` file and set `TELEGRAM_BOT_TOKEN` environment variable to be able to communicate with
 your bot via the Telegram bot API.
 * run `docker-compose build` to create the container image
 * run `docker-compose up -d` to create and start the container

## How to stop the web service

 * run `docker-compose down` to stop and destroy the container