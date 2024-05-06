FROM oven/bun:latest AS base

WORKDIR /app

COPY package.json bun.lockb ./
RUN bun install

# Copy the rest of the files
COPY . .

ARG BUILD_ENV=production

EXPOSE 3000

RUN if [ "$BUILD_ENV" = "production" ]; then \
        bun run build; \
    else \
        echo "Skipping build step"; \
    fi

CMD if [ "$BUILD_ENV" = "production" ]; then \
        echo "Starting production server"; \ 
        #todo: ADD PRODUCTION SERVER COMMAND
        bun run dev; \
    else \
        echo "Starting dev server"; \
        bun run dev; \
    fi