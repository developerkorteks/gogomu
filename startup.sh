#!/bin/bash

# Startup script for DOM Cloud Passenger
echo "🚀 Starting MultipleScrape API..."
echo "📊 Port: $PORT"
echo "🔧 Mode: $GIN_MODE"

# Set default values if not provided
export PORT=${PORT:-8080}
export GIN_MODE=${GIN_MODE:-release}

# Start the application
exec ./app