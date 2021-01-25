# Telegram Notification Service

This web service is able to send notifications to users via the Telegram bot API.

For being able to receive messages from a bot in Telegram a user has to start the initial communication by searching for the Bot's username and entering the `/start` command into the chat.
After this the bot should respond by providing a chat id which has to be set in the `recipient` field for being able to send further messages to this user.
To get more information about the endpoints of this service, please consult the OpenAPI documentation which is shipped via Swagger Editor.

There are currently two different message types sent by the bot. The first is the welcome message (`WELCOME_MESSAGE`) which shows when you
first start a conversation (by sending the bot command `/start`). The second type of message is called the unknown
message (`UNKNOWN_MESSAGE`) which is generated if you send anything other than `/start`. You are able to use this service in a little
more generic way by defining your own messages for each of both types.

In addition, there are some pre-defined templating variables maintaining information about the current user who is
chatting with the bot available. These are `USER_FIRST_NAME`, `USER_LAST_NAME`, `USER_FULL_NAME`, `USER_USERNAME` and
`USER_CHAT_ID`.
These variables can be included into the `WELCOME_MESSAGE` and `UNKNOWN_MESSAGE` message strings to let the responses appear being more personalized to the user.

## Environment Variables

The behavior of the application can be determined by setting some environment variables.

| VARIABLE NAME       | DEFAULT VALUE     | DESCRIPTION     |
| :------------- | :---------- | :----------- |
|  FLASK_APP | `app.py`  | Location of the Flask application    |
|  DB_FILE_PATH | `./db/db.sql`  | Location of the SQLite database file    |
|  TELEGRAM_BOT_TOKEN | `<empty>`  | Telegram Bot token   |
|  AUTH0_ISS | `scc2020g8.eu.auth0.com/`  | Token issuer. Required for Auth0 token based client credentials flow. |
|  AUTH0_AUD | `https://185.128.119.135/notification`  | Audience for the token. Required for Auth0 token based client credentials flow. |
|  WELCOME_MESSAGE | `Welcome, {{ USER_FULL_NAME }}! You are now able to configure {{ SERVICE_NAME }} via Telegram. Go to your dashboard and add the Telegram notification provider by entering this ID {{ USER_CHAT_ID }}.`  | Message the bot responds with after the user has started the conversation (bot command `/start`). |
|  UNKNOWN_MESSAGE | `Sorry, {{ USER_FULL_NAME }}. I didn't understand you.`  | Message the bot responds with after the user sent a message different from `/start`. |
|  TEMPLATE_VAR_SERVICE_NAME | `COVID-19 incidence update notifications` | Template variable which can be used in both message types `WELCOME_MESSAGE` and `UNKNOWN_MESSAGE`. Purpose of the application. |

## Prerequisites

1) In order to be able to send messages via Telegram you need create a
[Telegram bot](https://telegram.org/blog/bot-revolution).
Read [this](https://core.telegram.org/bots) introduction to get in touch with the Telegram bot platform.
You are able to generate a bot token after you created the bot.
2) Open [docker-compose.yml](docker-compose.yml) file and set `TELEGRAM_BOT_TOKEN` environment variable to be able to
   communicate with your bot via the Telegram bot API. It is also required to set `AUTH0_ISS` and `AUTH0_AUD`.
   For example `AUTH0_ISS=scc2020g8.eu.auth0.com/` and `AUTH0_AUD=https://185.128.119.135/notification`.
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
 * The service should be available at [http://localhost/notification](http://localhost/notification), and Swagger should listen at [http://localhost](http://localhost).

## How to stop the web service

 * Enter `docker-compose down` to stop and destroy the containers.

## Authorization

This web service uses [OAuth 2.0](https://tools.ietf.org/html/rfc6749) authorization via [Auth0](https://auth0.com).
Therefore, you have to set an `Authorization` header containing a bearer token in your requests. 
The token can be obtained and saved into the variable `AUTH_TOKEN` by executing the following request:
(Make sure you have both packages `curl` and the json parser `jq` installed on your system.)

```bash
$ CLIENT_SECRET=your-secret CLIENT_ID=your-client-id AUDIENCE=your-audience
$ AUTH_TOKEN=$(curl --request POST --url https://scc2020g8.eu.auth0.com/oauth/token \
                    --header 'content-type: application/json' \
                    --data "{\"client_id\":\"${CLIENT_ID}\", \
                             \"client_secret\":\"${CLIENT_SECRET}\", \
                             \"audience\":\"${AUDIENCE}\", \
                             \"grant_type\":\"client_credentials\"}" | jq -r .access_token)
```

Now you should be able to do requests:

```bash
$ curl localhost/notification -H "Accept: application/hal+json" -H "Authorization: Bearer ${AUTH_TOKEN}"
```

Detailed information about client credentials flow can be found
[here](https://auth0.com/docs/flows/call-your-api-using-the-client-credentials-flow).

## Running tests

In `tests` dir there is a number of API endpoint tests which can be performed by using
[Postman](https://www.postman.com/). After Postman is installed the test collection can be imported by using `File` and
`Import...` ** Since some tests depend on each other (by setting Postman collection wide variables) the tests should
run in the predefined order. In addition, it is also required to set a valid bearer token in the config section of the
collection. **