# Covid19 Update Service

This service allows subscribing to the 7 day Covid 19 incidence value provided by the "Robert Koch Institut" (RKI
) for a specific geo location.

RKI API: https://npgeo-corona-npgeo-de.hub.arcgis.com/datasets/917fc37a709542548cc3be077a786c17_0

## Setup

Use the docker-compose file to run the service locally.

### Environment Variables

|  Variable   | Description |
|-------------|-------------|
|   DB_TYPE   | Database type, currently supported: `sqlite3` and `mysql` |
|  DB_SOURCE  | Connection string, see: https://gorm.io/docs/connecting_to_the_database.html |
| SERVER_PORT | Port the server binds to. This port has to be exposed! |
| POLL_INTERVAL_MINUTES | Interval (in minutes) in which the data is retrieved from the RKI API |
| TELEGRAM_CONTACT_URI | URL of Telegram notification REST API |
| SENDGRID_API_KEY | API Key for SendGrid (https://sendgrid.com) |
| AUTH0_ISS | Auth0 issuer |
| AUTH0_AUD | Auth0 API ID |
| AUTH0_AUD | Auth0 Realm |