# Set Up for Development

## Config

todo: update config explanation

## Spinning up locally

```
BUILD_ENV=development docker-compose up --build
```

## Usage

After the containers have been set up locally, the client can be accessed at [http://localhost:3000]. The server can be accessed at [http://localhost:3001].

## Shutting down

```
docker-compose down --remove-orphans
```

## Troubleshooting

- For node_module errors:
  1. delete the node modules folder
  2. cd into client
  3. run `npm install`
  4. spin up containers

## For deployment

todo update section

```
BUILD_ENV=production docker-compose up --build
```

---

# Architectural Information

## Client

- language: javascript
- framework: Nuxt3
- architecture follows framework conventions
- Components:
  - components subfolders are arranged by function
  - "ui" folder contains shadcn-vue components
    - the tailwind classes of the components in the components folder are global
    - the tailwind classes of the components in the pages override the global if applied
- Store:
  - library: pinia, pinia-persisted-state

## Server

- language: Go
- server framework: Echo
- database: PostgreSQL
- ORM: [GORM] (https://gorm.io/)

### Style

- follow [Uber-Go Style Guide](https://github.com/uber-go/guide) wherever possible
- architecture is domain driven
- follow CLEAN & SOLID principles wherever feasible

### Architecture

- CLEAN architecture
- Partially Domain Driven
- dependency injection: hander (api interface) <- service (business logic) <- repository (db actions) <- db client
- todo: define narrower interfaces

### Folder Structure

All internal domains are under the internal folder.

Each domain contains:

- Model: structs specific to the domain
- Error: errors specific to the domain
- Routes: API routes specific to the domain
- Handler: API interfaces specific to the domain
  1. Claims validation & setting
  2. Cookies validation & setting
  3. Marshaling & Unmarshaling JSON
  4. Other validations
- Service: business logic specific to the domain
  1. Core business logic
  2. Ignorant of clientside interactions
  3. Ignorant of database interactions
- Repository: db actions specific to the domain
  1. Database actions

Common domain:

- Contains Models, Errors, etc. common between various domains

Schema domain:

- Contains database schema structs
- Reason why they are not under specific domains is due to cyclical imports. To refactor later.
  (I don't know an elegant solution to this problem. Would welcome advice.)

## Database

- A local postgreSQL database included in docker-compose setup
- startup/db.go will establish a connection on spin up, will panic if the connection fails
- if using pgAdmin4 to connect to the local postgres container, use `host:localhost`

---

# API Documentation

- library: [swaggo](https://github.com/swaggo/swag)
- middleware: [echo swagger](https://github.com/swaggo/echo-swagger)
- to access the swagger page:

```
http://localhost:3001/swagger
```

### Generate swagger documentation

```
cd server
swag init
```

### Accessing swagger documentation

[SWAGGER LINK](http://localhost:3001/swagger/index.html)

---

### Git Conventions

Use [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/#summary)
