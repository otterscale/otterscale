# Build stage
FROM node:25-alpine@sha256:b9b5737eabd423ba73b21fe2e82332c0656d571daf1ebf19b0f89d0dd0d3ca93 AS builder

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
FROM node:25-alpine@sha256:b9b5737eabd423ba73b21fe2e82332c0656d571daf1ebf19b0f89d0dd0d3ca93

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