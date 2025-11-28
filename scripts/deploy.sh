#!/bin/sh

set -e

echo "🔍 Checking environment files..."

# check if .env file exists in apps/dashboard
if [ ! -f ./apps/dashboard/.env ]; then
  echo "⚠️  .env file not found in apps/dashboard. Creating default .env file"
  echo "VITE_API_BASE_URL=https://hc-ppn-chatbot.cloud/api" > ./apps/dashboard/.env
  echo "✅ .env file created"
else
  echo "✅ .env file exists"
fi

# check if .env file exists in apps/bot-service/config
if [ ! -f ./apps/bot-service/config/.env ]; then
  echo "⚠️  .env file not found in apps/bot-service/config. Copying from .env.example"
  cp ./apps/bot-service/config/.env.example ./apps/bot-service/config/.env
  echo "✅ .env file created"
else
  echo "✅ .env file exists"
fi

echo ""
echo "🐳 Building and starting Docker containers..."
docker compose -f ./docker/docker-compose.yaml up -d --build

echo ""
echo "⏳ Waiting for database to be ready..."
sleep 5

# Wait for database to be healthy
until docker exec hc-ppn-postgres pg_isready -U postgres -d hc_ppn_app > /dev/null 2>&1; do
  echo "   Database is not ready yet, waiting..."
  sleep 2
done

echo "✅ Database is ready"

echo ""
echo "🔄 Running database migrations..."

# Load .env file to get database credentials
export $(cat ./apps/bot-service/config/.env | grep -v '^#' | sed 's/#.*$//' | grep -v '^[[:space:]]*$' | xargs)

# Construct database URL
DB_URL="postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

# Run migrations using docker with traefik network
docker run --rm \
  -v ./apps/bot-service/database/migrations:/database/migrations \
  --network docker_hc-ppn-network \
  migrate/migrate \
  -path /database/migrations \
  -database "$DB_URL" \
  -verbose up

echo ""
echo "✅ Deployment completed successfully!"
echo ""
echo "📊 Services status:"
docker ps
