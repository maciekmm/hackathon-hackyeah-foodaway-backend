# HackYeah Entry

Foodsharing idea implementation.

## Running

In order to start the backend you need to have a docker-compose and run the following commands:

```bash
cp ./credentials.env.sample ./credentials.env
vi ./credentials.env #change values in credentials.env to yours
docker-compose up --build
```

### Generating fake data

```bash
docker exec hackyeah_api_1 go run cmd/mock/main.go N
```

where N is a number of records to be generated

## Endpoints

### POST /accounts/register/

Sample request:

```json
{
    "name": "Your Name",
    "email": "email@example.com",
    "password": "your-password"
}
```

Sample response:

```json
{
    "token": "JWT-token"
}
```

All future requests have to contain `Authorization: Bearer YOUR-JWT-TOKEN` header. Token is valid for 24h

### POST /accounts/login/

```json
{
    "email": "email@example.com",
    "password": "your-password"
}
```

Sample response:

```json
{
    "token": "JWT-token"
}
```

### GET /accounts/token/

Used to refresh expired token. Pass old token in `Authorization` header and a new one will be returned.

Sample response:

```json
{
    "token": "new-jwt-token"
}
```

