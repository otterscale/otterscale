# Build stage
FROM node:24-alpine@sha256:c921b97d4b74f51744057454b306b418cf693865e73b8100559189605f6955b8 AS builder

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
FROM node:24-alpine@sha256:c921b97d4b74f51744057454b306b418cf693865e73b8100559189605f6955b8

WORKDIR /app

# Install pnpm
RUN npm install -g pnpm

# Copy package.json and pnpm-lock.yaml
COPY package.json pnpm-lock.yaml ./

# Install only production dependencies
RUN pnpm install --production --frozen-lockfile && pnpm store prune

# Copy built application from builder stage
COPY --from=builder --chown=nobody:nobody /src/build ./build

# Switch to non-root user
USER nobody

# Set environment variable
ENV NODE_ENV=production

# Expose port
EXPOSE 3000

# Labels
LABEL maintainer="Chung-Hsuan Tsai <paul_tsai@phison.com>"

CMD ["node", "./build"]