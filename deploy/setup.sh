#!/bin/bash

# Fleming Deployment Script for Lightsail
# Usage: ./setup.sh

set -e

echo "🚀 Starting Fleming Deployment..."

# 1. Check for .env file
if [ ! -f .env ]; then
    echo "❌ Error: .env file not found!"
    echo "PLEASE CREATE .env FILE WITH THE FOLLOWING VARIABLES:"
    echo "DOMAIN_NAME=x.x.x.x.nip.io"
    echo "ACME_EMAIL=your-email@example.com"
    echo "DATABASE_URL=postgres://..."
    echo "JWT_SECRET=..."
    exit 1
fi

# 2. Determine Docker Compose command
if docker compose version >/dev/null 2>&1; then
    COMPOSE="docker compose"
elif docker-compose version >/dev/null 2>&1; then
    COMPOSE="docker-compose"
else
    echo "❌ Error: Docker Compose not found!"
    echo "Please install it: sudo curl -L \"https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)\" -o /usr/local/bin/docker-compose && sudo chmod +x /usr/local/bin/docker-compose"
    exit 1
fi

# 3. Pull latest images
echo "📥 Pulling latest images..."
$COMPOSE -f compose.prod.yml pull

# 4. Restart services
echo "🔄 Restarting services..."
$COMPOSE -f compose.prod.yml up -d --remove-orphans

echo "✅ Deployment Complete!"
echo "---------------------------------------------------"
echo "Backend is running at: https://$(grep DOMAIN_NAME .env | cut -d '=' -f2)"
