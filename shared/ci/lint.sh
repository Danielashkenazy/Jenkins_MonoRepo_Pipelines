#!/bin/bash
set -e

SERVICES_CHANGED=${1:-$SERVICES_CHANGED}

if [ -z "$SERVICES_CHANGED" ]; then
    echo "No services to lint"
    exit 0
fi

echo "Services to lint: $SERVICES_CHANGED"

for SERVICE in $SERVICES_CHANGED; do
    echo "===================="
    echo "Linting: $SERVICE"
    echo "===================="
    
    cd "$SERVICE"
    
    if [ -f "package.json" ]; then
        echo "Running ESLint for Node.js service..."
        npm install
        npx eslint .
        
    elif [ -f "requirements.txt" ]; then
        echo "Running Flake8 for Python service..."
        python3 -m venv .venv
        . .venv/bin/activate
        pip install -r requirements.txt
        pip install flake8
        flake8 . --exclude=.venv,__pycache__,.git
        
    elif [ -f "go.mod" ]; then
        echo "Running golangci-lint for Go service..."
        go mod tidy
        go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
        export PATH=$(go env GOPATH)/bin:$PATH
        golangci-lint run
    else
        echo "Unknown service type for $SERVICE"
        exit 1
    fi
    
    cd ..
    echo "✅ Linting completed for $SERVICE"
done