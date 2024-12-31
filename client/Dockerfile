# Base stage
FROM oven/bun:latest AS base

WORKDIR /app

# Copy package.json and bun.lockb first for better cache management
COPY package.json bun.lockb ./

# Install dependencies using Bun
RUN bun install

# Copy the rest of the files
COPY . .

# Define build arguments
ARG BUILD_ENV=production

EXPOSE 3000

# Build the project based on the environment
RUN if [ "$BUILD_ENV" = "production" ]; then \
        echo "Building for production"; \
        bun run build; \
    else \
        echo "Skipping build step for development"; \
    fi

# Define the command to run based on the environment
CMD if [ "$BUILD_ENV" = "production" ]; then \
        echo "Starting production server"; \
        bun .output/server/index.mjs; \
    else \
        echo "Starting dev server"; \
        bun run dev; \
    fi
