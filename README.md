# mfv-challenge

## Install

This application is using Go 1.21 and Postgres as a database. To run this application, should pull this repository from Github and run by docker compose like this
```sh
docker-compose -f docker-compose.yml up
```
The database schema will be run by the migration in the docker compose with the mock data for users and accounts as file [mock data](./migrations/000002_mock_data.up.sql)

To set the application mannually, pull the depedencies by command `make dep`, install and run Postgres. After that, install the [go-migrate](https://github.com/golang-migrate/migrate) to update the database schema.

Then running the one of 2 commands to migrate the schema & mock data

```sh
$ migrate -source file://path/to/migrations -database postgres://localhost:5432/database up 3
```

```sh
$ docker run -v {{ migration dir }}:/migrations --network host migrate/migrate
    -path=/migrations/ -database postgres://localhost:5432/database up 3
```

The migration files are on [this directory](./migrations/)

## APIs

As the requirements, it has 3 APIs that follows the coding challenge as below

### List all accounts of an user
```sh
curl --location 'localhost:8080/api/users/1/accounts'

[
    {
        "id": 1,
        "user_id": 1,
        "name": "alice 1",
        "balance": 12000
    },
    {
        "id": 3,
        "user_id": 1,
        "name": "alice 2",
        "balance": 12000
    },
    {
        "id": 5,
        "user_id": 1,
        "name": "alice 3",
        "balance": 12000
    }
]
```

### Get an user and list account's id

```sh
curl --location 'localhost:8080/api/users/1'

{
    "id": 1,
    "name": "alice",
    "account_ids": [
        1,
        3,
        5
    ]
}
```

### 

```sh
curl --location 'localhost:8080/api/accounts/1'

{
    "id": 1,
    "user_id": 1,
    "name": "alice 1",
    "balance": 12000
}
```

## Supported APIs

### List transactions

```sh
curl --location 'localhost:8080/api/users/1/transactions'

[
    {
        "id": 3,
        "account_id": 3,
        "amount": 40000,
        "transaction_type": "deposit",
        "created_at": "2024-01-24T10:14:52.905879Z"
    },
    {
        "id": 2,
        "account_id": 3,
        "amount": 50000,
        "transaction_type": "deposit",
        "created_at": "2024-01-24T10:14:48.885647Z"
    },
    {
        "id": 1,
        "account_id": 3,
        "amount": 100000,
        "transaction_type": "deposit",
        "created_at": "2024-01-24T10:14:22.786115Z"
    }
]
```

### Create a transaction

```sh
curl --location 'localhost:8080/api/users/1/transactions' \
--header 'Content-Type: application/json' \
--data '{
"account_id": 3,
"amount": 100000.00,
"transaction_type": "deposit"
}'

{
    "id": 1,
    "account_id": 3,
    "amount": 100000,
    "bank": "ACB",
    "transaction_type": "deposit",
    "created_at": "2024-01-24T10:14:22.786115795Z"
}
```

## Authentication

This application supports a simple mechanism for authentication with JWT. To enable the authentication method, should set the variable `enabled` of `auth` as true (default is false) in the [config file](./config/config.yml) or update environment of `api` in the docker-compose file

To get the token, use `login` API:

```sh
curl --location 'localhost:8080/api/login' \
--header 'Content-Type: application/json' \
--data '{
    "username": "alice",
    "password": "123456"
}'

{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3MDYwOTQ1NjEsIm5iZiI6MTcwNjA5MDk2MSwiaWF0IjoxNzA2MDkwOTYxfQ.E2oZCFTuQE836PYjSBl_Geh3FFRaf5fh_gaigsWlmC8"
}
```

Add the Authorization Header to each API request when enabled the authentication as below:

```sh
curl --location 'localhost:8080/api/accounts/1' \
--header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJleHAiOjE3MDYwNzk3MDcsIm5iZiI6MTcwNjA3NjEwNywiaWF0IjoxNzA2MDc2MTA3fQ.MO2CCcXt4mbIFR4lyoGgJKs81qam6Uqn5edgak7bCz4'

```