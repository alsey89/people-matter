# Set Up for Development

## Local Environmental Variables

A config.yaml file with sensible default/temporary values **is included** in the repository. This is for convenience during development. This should **NOT** be used in production.

Config files are managed with:

- viper: https://github.com/spf13/viper

Default behavior:

- if environmental variables are supplied, viper will use **environmental variables**
- else viper will fall back to the default **config.yaml**
- if **config.override.yaml** exists, it will override config.yaml

## Spinning up locally

```
docker-compose up --build
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
- For "cannot find defineNuxtConfig" errors
  - try:
    1. run `npx nuxi cleanup`
    2. step 1 to 4 under "for node_module errors"
  - if problem persists, try: 3. setting up [volar takeover](https://vuejs.org/guide/typescript/overview#volar-takeover-mode) if using VS Code

## For deployment

- edit docker-compose.deploy.yaml
- BASE_API_URL should point to the backend directly (OR)
- BASE_API_URL can point to the domain and Nginx can route it to the backend

- when deploying as a docker-compose setup use:

```
docker-compose -f docker-compose.yaml -f docker-compose.deploy.yaml up -d
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
