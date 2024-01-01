# Dev Environment

## Set up using docker-compose

### Spinning up

```
docker-compose build
docker-compose up
```

### Shutting down

```
docker-compose down --remove-orphans
```

## troubleshooting

- for node_module errors try deleting the node modules folder, then spin up docker-compose
- for "cannot find defineNuxtConfig" errors, try:
  - setting up [volar takeover](https://vuejs.org/guide/typescript/overview#volar-takeover-mode) if using VS Code
  - deleting the node modules folder and the .nuxt folder. Then, do an npm install outside of the container (client) before spinning up docker-compose.

# Architecture

## Client is Nuxt3

- architecture follows framework conventions
- components subfolders are arranged by function
- NB folder contains neobrutalism design components

## Server is Go-Echo

- style should follow [Uber-Go Style Guide](https://github.com/uber-go/guide) wherever possible
- architecture is domain driven
- follow CLEAN & SOLID principles wherever feasible
- dependency injection: hander(api interface) <- service (business logic) <- repository (db)
- todo: define narrow interfaces
