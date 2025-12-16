# Build stage
FROM node:24-alpine AS builder

ARG VERSION=unknown

WORKDIR /src

# Install pnpm
RUN npm install -g pnpm

# Copy dependency files first for better caching
COPY package.json pnpm-lock.yaml ./
RUN npm install -g pnpm && pnpm install --frozen-lockfile && pnpm store prune

# Copy source code
COPY . .

# Set version environment variable
ENV VERSION=${VERSION}

# Build the application
RUN pnpm build

# Runtime stage
FROM node:24-alpine

WORKDIR /app

# Install pnpm
RUN npm install -g pnpm

# Copy package.json and pnpm-lock.yaml
COPY package.json pnpm-lock.yaml ./

# Install only production dependencies
RUN pnpm install --production --frozen-lockfile && pnpm store prune

# Copy built application from builder stage
COPY --from=builder --chown=nobody:nobody /src/build ./build

# Copy migration script and drizzle migration files
COPY --from=builder --chown=nobody:nobody /src/scripts/migrate.js ./scripts/migrate.js
COPY --from=builder --chown=nobody:nobody /src/drizzle ./drizzle

# Switch to non-root user
USER nobody

# Set environment variable
ENV NODE_ENV=production

# Expose port
EXPOSE 3000

# Labels
LABEL maintainer="Chung-Hsuan Tsai <paul_tsai@phison.com>"

CMD ["/bin/sh", "-c", "node ./scripts/migrate.js && node ./build/index.js"]