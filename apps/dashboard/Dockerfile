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
FROM nginx:alpine

# Copy built files to nginx html directory
COPY --from=builder /app/apps/dashboard/dist /usr/share/nginx/html

# Copy nginx configuration
COPY apps/dashboard/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
