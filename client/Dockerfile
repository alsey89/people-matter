# Base stage
FROM oven/bun:latest AS base

WORKDIR /app

COPY package.json bun.lockb ./
RUN bun install

# Copy the rest of the files
COPY . .

# Define build arguments
ARG BUILD_ENV=production

EXPOSE 3000

RUN if [ "$BUILD_ENV" = "production" ]; then \
        echo "Building for production"; \
        bun run build; \
    else \
        echo "Skipping build step for development"; \
    fi

CMD if [ "$BUILD_ENV" = "production" ]; then \
        echo "Starting production server"; \
        # todo: add production server start command
        # bun run start; \
    else \
        echo "Starting dev server"; \
        bun run dev; \
    fi
