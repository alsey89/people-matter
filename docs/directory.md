## Directory Structure

The following is the directory structure of the project:

```bash
├── client/                     # Frontend code
├── database/                   # Database setup
├── docs/                       # Documentation
├── server/                     # Backend code
├── .gitignore
├── docker-compose.yaml
└── README.md
```

Client directory structure

```bash
client/
├── public/                # Static assets
│   ├── *
├── src/
│   ├── assets/            # Assets like images, fonts, etc.
│   │   └── *
│   ├── components/
│   │   └── *
│   ├── layouts/
│   │   └── *
│   ├── pages/
│   │   ├── **/*.vue
│   ├── plugins/
│   │   ├── *.ts
│   ├── router/
│   │   ├── index.ts       # Router setup
│   │   └── routes.ts      # Route definitions
│   ├── stores/            # Pinia store
│   │   ├── *.ts
│   ├── App.vue            # Root component
│   ├── main.ts            # Entry point for the application
│   └── utils/             # Utility functions and helpers
│       └── *
├── .env.dev
├── .env.prod
├── Dockerfile
├── package.json
├── .env.dev
├── *
```

Server directory structure

```bash
server/
├── config/
│   ├── *.go               # Configuration files
├── internal/
│   ├── domain/
│   │   ├── */domain.go    # Domain struct & logic
│   │   ├── */error.go     # Custom error handling
│   │   ├── */handler.go   # HTTP handler
│   │   ├── */model.go     # Non-Schema Data models
│   │   ├── */service.go   # Service layer
│   │   ├── */test.go      # Domain struct & logic
│   ├── schema/
│   │   ├── *.go           # Schema definitions
├── pkg/
│   ├── **/*.go            # Reusable packages and utilities
├── Dockerfile
├── go.mod
├── go.sum
└── main.go                # Main entry point of the application
```
