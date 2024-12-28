# Base image
FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV GOBIN=/usr/local/bin
ARG BUILD_ENV=production
ENV SERVER_BUILD_ENV=$BUILD_ENV

EXPOSE 3001

RUN if [ "$BUILD_ENV" = "development" ]; then \
        go install github.com/air-verse/air@latest; \
    fi

CMD if [ "$BUILD_ENV" = "development" ]; then \
        echo "Starting dev server"; \
        air; \
    else \
        echo "Starting production server"; \
        go build -o /bin/main . && /bin/main; \
    fi
