# Dev Environment

## Docker-compose for setup

### Spinning up

```
docker-compose build
docker-compose up
```

### Usage

After the containers have been set up, the application FE can be accessed at [http://localhost:3000]. The server can be accessed at [http://localhost:3001]. Since the project is in development, the env file is included in the repo and secure features have been disabled.

### Shutting down

```
docker-compose down --remove-orphans
```

### Troubleshooting

- for node_module errors try deleting the node modules folder, then spin up docker-compose
- for "cannot find defineNuxtConfig" errors, try:
  - setting up [volar takeover](https://vuejs.org/guide/typescript/overview#volar-takeover-mode) if using VS Code
  - deleting the node modules folder and the .nuxt folder. Then, do an npm install outside of the container (client) before spinning up docker-compose.

---

# Documentation

## Server

- library: [swaggo](github.com/swaggo/swag)
- middleware: [echo swagger](https://github.com/swaggo/echo-swagger)

### Generate swagger documentation

```
cd server
swag init
```

### Accessing swagger documentation

[SWAGGER LINK](http://localhost:3001/swagger/index.html)

---

# Architecture

## Client is Nuxt3

- architecture follows framework conventions
- components subfolders are arranged by function
- NB folder contains neobrutalism design components

## Server

- language: Go
- server framework: Echo
- database: PostgreSQL
- ORM: [GORM] (https://gorm.io/)

### Style

- follow [Uber-Go Style Guide](https://github.com/uber-go/guide) wherever possible
- architecture is domain driven
- follow CLEAN & SOLID principles wherever feasible
- dependency injection: hander(api interface) <- service (business logic) <- repository (db)
- todo: define narrow interfaces

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
- start/db.go will establish a connection on spin up, will panic if the connection fails
- if using pgAdmin4 to connect to the local postgres container, use `host:localhost`
