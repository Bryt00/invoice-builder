#!/bin/bash
echo "========================================"
echo " Starting Frontend-Only Update Script "
echo "========================================"

# Resolve the directory of the script and navigate to the project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_ROOT"
echo "Working directory set to: $PROJECT_ROOT"

APP_DIR="/opt/invoice-builder"

if [ ! -d "$APP_DIR" ]; then
    echo "Error: Application directory $APP_DIR does not exist. Run deploy.sh first."
    exit 1
fi

echo "Building Vue frontend..."
if [ -d "frontend" ]; then
    cd frontend
    
    # Optional: You can uncomment this if you added new npm packages
    # npm cache clean --force
    # npm install --no-audit --no-fund
    
    npm run build
    cd ..
    
    echo "Deploying built files to $APP_DIR/frontend/dist..."
    mkdir -p $APP_DIR/frontend/dist
    cp -r frontend/dist/* $APP_DIR/frontend/dist/
    
    echo "========================================"
    echo " Frontend successfully updated! "
    echo "========================================"
else
    echo "Error: 'frontend' directory not found."
    exit 1
fi
