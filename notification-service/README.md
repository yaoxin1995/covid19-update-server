# Notification Service

This web service is able to send notifications to users via the Telegram bot API.

There are currently two different message types sent by the bot. The first is the welcome message which shows when you
first start a conversation (by sending the bot command `/start`). The second type of message is called the unknown
message which is generated if you send anything other than `/start`. You are able to use this service in a little
more generic way by defining you own messages for each of both types.

This project is using [Jinja templating engine](https://palletsprojects.com/p/jinja/) for rendering messages with
variables. Have a look at [docker-compose.yml](docker-compose.yml) file to find out how you can use environment
variables in messages.
In addition, there are some pre-defined variables available maintaining information about the current user who is
chatting with the bot. These are `USER_FIRST_NAME`, `USER_LAST_NAME`, `USER_FULL_NAME`, `USER_USERNAME` and
`USER_CHAT_ID`. The latter one is needed if you want to directly send messages to a user via the bot.
Have a look at the API documentation for more details.

## Prerequisites

1) In order to be able to send messages via Telegram you need create a
[Telegram bot](https://telegram.org/blog/bot-revolution).
Read [this](https://core.telegram.org/bots) introduction to get in touch with the Telegram bot platform.
You are able to generate a bot token after you created the bot.
2) Open [docker-compose.yml](docker-compose.yml) file and set `TELEGRAM_BOT_TOKEN` environment variable to be able to
   communicate with your bot via the Telegram bot API.
3) By default, this web service uses HTTPS. Therefore, you need to create a certificate and corresponding private key.
   This can be done by this command (if you don't own a certificate already):
   `mkdir -p ./nginx/cert && openssl req -x509 -newkey rsa:4096 -nodes -out ./nginx/cert/cert.pem -keyout ./nginx/cert/key.pem -days 365`.
   This directory has to be bind-mounted into the container which is done in the `volumes` section of the Docker compose
   file.
4) *optional:* Set your own messages (`WELCOME_MESSAGE` and `UNKNOWN_MESSAGE`) via environment variables in
   [docker-compose.yml](docker-compose.yml) file.
5) Run `docker-compose build` to create the container image.

## How to run the web service

By default, three containers are created: the web service itself, Swagger Editor, and Nginx which is used as a reverse
proxy. To prevent trouble with CORS, Swagger is also behind the reverse proxy.

 * Enter `docker-compose up -d` to create and start the containers.

## How to stop the web service

 * Enter `docker-compose down` to stop and destroy the container.