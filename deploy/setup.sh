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

# 2. Pull latest images
echo "📥 Pulling latest images..."
docker compose -f compose.prod.yml pull

# 3. Restart services
echo "🔄 Restarting services..."
docker compose -f compose.prod.yml up -d --remove-orphans

echo "✅ Deployment Complete!"
echo "---------------------------------------------------"
echo "Backend is running at: https://$(grep DOMAIN_NAME .env | cut -d '=' -f2)"
