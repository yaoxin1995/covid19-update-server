# Covid19 Update Service

This service allows subscribing to the COVID-19 7-Day-Incidence value provided by the 
"Robert Koch Institut" (RKI) for a chosen geo location. A threshold value can be configured for each location. Once 
the threshold is exceeded a notification will be sent.
Currently Telegram and email notifications are supported.

## Usage

See `/doc` for the OpenAPI specification.

### Authorization

To use the service, authorization has to be performed via OAuth 2.0 using Auth0 (https://auth0.com).

The OAuth 2.0 access token returned from auth0 is a JWT token that has to be used as Bearer token inside the 
authorization header.

Access will be restricted to resources, that are owned by the subject referred to in the token's `sub` claim.

#### Client Credential Flow

To use the client credential flow (https://tools.ietf.org/html/rfc6749#section-4.4) follow the instructions of Auth0 
(https://auth0.com/docs/flows/call-your-api-using-the-client-credentials-flow), e.g.:

```
curl --request POST \
  --url https://scc2020g8.eu.auth0.com/oauth/token \
  --header 'content-type: application/json' \
  --data '{"client_id":"CLIENT_ID","client_secret":"CLIENT_SECRET",
  "audience":"https://185.128.119.135","grant_type":"client_credentials"}'
```

(Replace `ClIENT_ID` and `ClIENT_SECRET` with the ID and secret of for your application)

## Dependencies

- RKI API: https://npgeo-corona-npgeo-de.hub.arcgis.com/datasets/917fc37a709542548cc3be077a786c17_0
- SendGrid (for email notifications): https://sendgrid.com
- Telegram Notification Service
- Auth0 (for client authentication): https://auth0.com

## Setup

Use the docker-compose file to run the service locally.

### Environment Variables

|  Variable   | Description |
|-------------|-------------|
|   DB_TYPE   | Database type, currently supported: `sqlite3` and `mysql` |
|  DB_SOURCE  | Connection string, see: https://gorm.io/docs/connecting_to_the_database.html |
| SERVER_PORT | Port the server binds to. This port has to be exposed! |
| CORS_ORIGINS| Comma separated list of allowed origins for CORS |
| POLL_INTERVAL_MINUTES | Interval (in minutes) in which the data is retrieved from the RKI API |
| TELEGRAM_CONTACT_URI | URL of Telegram notification REST API |
| SENDGRID_API_KEY | API Key for SendGrid |
|SENDGRID_EMAIL| Email address used to send email notifications |
| AUTH0_ISS | Auth0 issuer, e.g. `https://scc2020g8.eu.auth0.com/` (Trailing slash required!) |
| AUTH0_AUD | Auth0 API ID, e.g. `https://185.128.119.135` |

## Run Tests

To execute the test collection in `/test` you have to install Postman: https://www.postman.com

1. Import and select environment `/test/COVID-19 Update Service Test-Env.postman_environment.json`
2. Import and execute test collection `/test/COVID-19 Update Service.postman_collection.json`