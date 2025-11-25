#!/bin/bash
set -e

SERVICES_CHANGED=${1:-$SERVICES_CHANGED}

if [ -z "$SERVICES_CHANGED" ]; then
    echo "No services to test"
    exit 0
fi

echo "Services to test: $SERVICES_CHANGED"

for SERVICE in $SERVICES_CHANGED; do
    echo "===================="
    echo "Testing: $SERVICE"
    echo "===================="
    
    cd "$SERVICE"
    
    if [ -f "package.json" ]; then
        echo "Running tests for Node.js service..."
        npm install
        npm test -- --coverage --reporters=default --reporters=jest-junit
        
    elif [ -f "requirements.txt" ]; then
        echo "Running tests for Python service..."
        python3 -m venv .venv || true
        . .venv/bin/activate
        pip install pytest pytest-cov
        pytest --cov=. --cov-report=html --cov-report=xml --junitxml=junit.xml
        
    elif [ -f "go.mod" ]; then
        echo "Running tests for Go service..."
        go mod tidy
        go test -v -coverprofile=coverage.out -covermode=atomic ./...
        go tool cover -html=coverage.out -o coverage.html
    else
        echo "Unknown service type for $SERVICE"
        exit 1
    fi
    
    cd ..
    echo "✅ Tests completed for $SERVICE"
done