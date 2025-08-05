#!/bin/bash

# Build script for MultipleScrape API
echo "🚀 Building MultipleScrape API..."

# Clean previous build
echo "🧹 Cleaning previous build..."
rm -f app

# Download dependencies
echo "📦 Downloading dependencies..."
go mod download

# Verify dependencies
echo "✅ Verifying dependencies..."
go mod verify

# Build application
echo "🔨 Building application..."
go build -ldflags="-s -w" -o app

# Check if build successful
if [ -f "app" ]; then
    echo "✅ Build successful!"
    echo "📊 Binary size: $(du -h app | cut -f1)"
    echo "🎯 Ready for deployment!"
else
    echo "❌ Build failed!"
    exit 1
fi

# Test run (optional)
if [ "$1" = "test" ]; then
    echo "🧪 Testing application..."
    PORT=8080 ./app &
    APP_PID=$!
    
    sleep 3
    
    echo "🔍 Testing health endpoint..."
    curl -s http://127.0.0.1:8080/health
    
    echo ""
    echo "🔍 Testing monitoring endpoint..."
    curl -s http://127.0.0.1:8080/monitoring | head -c 200
    
    echo ""
    echo "🛑 Stopping test server..."
    kill $APP_PID
fi

echo "🎉 Build process completed!"