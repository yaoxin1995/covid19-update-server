# Covid19 Update Service

## Setup

Use the docker-compose file to run the service locally.

### Environment Variables

|  Variable   | Description |
|-------------|-------------|
|   DB_TYPE   | Database type, currently supported: `sqlite3` and `mysql` |
|  DB_SOURCE  | Connection string, see: https://gorm.io/docs/connecting_to_the_database.html |
| SERVER_PORT | Port the server binds to. This port has to be exposed! |
