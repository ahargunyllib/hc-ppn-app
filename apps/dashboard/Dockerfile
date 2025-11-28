# Build stage
FROM node:20-alpine as builder

WORKDIR /app

# Copy workspace configuration and lockfile
COPY pnpm-workspace.yaml pnpm-lock.yaml package.json ./

# Enable corepack for pnpm
RUN corepack enable

# Copy dashboard app package.json
COPY apps/dashboard/package.json apps/dashboard/

# Install dependencies (will install all workspace dependencies based on lockfile)
RUN pnpm install --frozen-lockfile

# Copy dashboard source code
COPY apps/dashboard/ apps/dashboard/

# Build dashboard
RUN pnpm --filter dashboard build

# Production stage
FROM oven/bun:1.3-alpine

WORKDIR /app

# Copy built files from builder
COPY --from=builder /app/apps/dashboard/dist .

EXPOSE 3001

CMD ["bun", "serve", "-p", "3001", "."]
