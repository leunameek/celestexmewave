#!/bin/bash

# CelestexMewave Frontend Start Script
# This script starts a local development server for the frontend

set -e

echo "🚀 Starting CelestexMewave Frontend..."
echo ""

# Check if Python 3 is installed
if command -v python3 &> /dev/null; then
    echo "✓ Python 3 found"
    echo ""
    echo "🔧 Starting frontend server on http://localhost:3000"
    echo ""
    python3 -m http.server 3000
    exit 0
fi

# Check if Python is installed
if command -v python &> /dev/null; then
    echo "✓ Python found"
    echo ""
    echo "🔧 Starting frontend server on http://localhost:3000"
    echo ""
    python -m http.server 3000
    exit 0
fi

# Check if Node.js is installed
if command -v node &> /dev/null; then
    echo "✓ Node.js found"
    echo ""
    
    # Check if http-server is installed globally
    if command -v http-server &> /dev/null; then
        echo "✓ http-server found"
        echo ""
        echo "🔧 Starting frontend server on http://localhost:3000"
        echo ""
        http-server -p 3000
        exit 0
    else
        echo "📦 Installing http-server..."
        npm install -g http-server
        echo ""
        echo "🔧 Starting frontend server on http://localhost:3000"
        echo ""
        http-server -p 3000
        exit 0
    fi
fi

# If we get here, no suitable server was found
echo "❌ Error: No suitable server found!"
echo ""
echo "Please install one of the following:"
echo "  • Python 3 (https://www.python.org/downloads/)"
echo "  • Node.js (https://nodejs.org/)"
echo ""
echo "Then run this script again."
exit 1
