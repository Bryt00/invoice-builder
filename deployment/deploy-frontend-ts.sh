#!/bin/bash
echo "========================================"
echo " Starting TypeScript Frontend Update Script "
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

echo "Building Vue TypeScript frontend..."
if [ -d "frontend_ts" ]; then
    cd frontend_ts
    
    if [ ! -d "node_modules" ]; then
        echo "Installing node_modules for frontend_ts..."
        npm install
    fi
    
    npm run build
    cd ..
    
    echo "Deploying built files to $APP_DIR/frontend/dist..."
    mkdir -p $APP_DIR/frontend/dist
    cp -r frontend_ts/dist/* $APP_DIR/frontend/dist/
    
    echo "========================================"
    echo " TypeScript Frontend successfully updated! "
    echo "========================================"
else
    echo "Error: 'frontend_ts' directory not found."
    exit 1
fi
