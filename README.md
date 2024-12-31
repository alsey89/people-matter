# CURATE MONOLITH

## Repository Structure

```bash
monolith/
├── client/
│   ├── Dockerfile
│
├── database/
│   ├── Dockerfile
│
├── server/
│   ├── Dockerfile
│   ├── config.yaml
│
├── docker-compose.yaml
└── README.md
```

## Dev Environment

Use the following command at root to spin up all containers in one command. Using the development BUILD_ENV will make both client and server hot-reload. Keep in mind that, occasionally, some errors will require you to do a full spin down and up.

```bash
BUILD_ENV=development docker-compose up

```

Use the following command at root to clean up all containers in one command.

```bash
docker-compose down --remove-orphans

```

## [server] Documentation

OpenAPI specifications and swagger are used for documentation. API handlers should contain detailed annotations upon opening pull requests.

Swagger UI can be found in the following location, assuming local development:
http://localhost:3001/swagger/index.html#/

Use the following command to generate the Swagger docs:

```
swag init --parseDependency --parseInternal
```
