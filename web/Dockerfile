# Build stage
FROM node:22-alpine AS builder

WORKDIR /src

# Install pnpm
RUN npm install -g pnpm

# Copy dependency files first for better caching
COPY package.json pnpm-lock.yaml ./
RUN npm install -g pnpm && pnpm install --frozen-lockfile && pnpm store prune

# Copy source code
COPY . .

# Build the application
RUN pnpm build

# Runtime stage
FROM node:22-alpine

WORKDIR /app

# Install pnpm
RUN npm install -g pnpm

# Copy package.json and pnpm-lock.yaml
COPY package.json pnpm-lock.yaml ./

# Install only production dependencies
RUN pnpm install --production --frozen-lockfile && pnpm store prune

# # Copy built application from builder stage
COPY --from=builder --chown=nobody:nobody /src/build ./build

# Switch to non-root user
USER nobody

# Set environment variable
ENV NODE_ENV=production

# Expose port
EXPOSE 3000

# Labels
LABEL maintainer="Chung-Hsuan Tsai <zx.c@nycu.edu.tw>"

CMD ["node", "./build"]