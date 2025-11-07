#!/bin/bash
# Script to restart the Bible Reading Backend server

echo "🛑 Stopping existing server..."
pkill -f "go run main.go" || pkill -f "bible_server" || echo "No server process found"

echo "⏳ Waiting for server to stop..."
sleep 2

echo "🔨 Building application..."
cd "$(dirname "$0")"
go build -o bible_server . || {
    echo "❌ Build failed!"
    exit 1
}

echo "✅ Build successful!"
echo ""
echo "🚀 Starting server..."
echo "   Server will run on http://localhost:8000"
echo "   Press Ctrl+C to stop"
echo ""

./bible_server

