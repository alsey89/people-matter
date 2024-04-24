# Base image
FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GOBIN=/usr/local/bin
ARG BUILD_ENV=production

EXPOSE 3001

RUN if [ "$BUILD_ENV" = "development" ]; then \
        go install github.com/cosmtrek/air@latest; \
    fi

CMD if [ "$BUILD_ENV" = "production" ]; then \
        echo "Starting production server"; \
        go build -o /bin/main . && /bin/main; \
    else \
        echo "Starting dev server"; \
        air; \
    fi
