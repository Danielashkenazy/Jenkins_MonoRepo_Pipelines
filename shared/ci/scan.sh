#!/bin/bash
set -e

SERVICES_CHANGED=${1:-$SERVICES_CHANGED}

if [ -z "$SERVICES_CHANGED" ]; then
    echo "No services to scan"
    exit 0
fi

echo "Services to scan: $SERVICES_CHANGED"

for SERVICE in $SERVICES_CHANGED; do
    echo "===================="
    echo "Security Scan: $SERVICE"
    echo "===================="
    
    cd "$SERVICE"
    
    if [ -f "package.json" ]; then
        echo "Running npm audit..."
        npm audit --audit-level=moderate
        
    elif [ -f "requirements.txt" ]; then
        echo "Running bandit for Python..."
        if [ ! -d ".venv" ]; then python3 -m venv .venv; fi
        . .venv/bin/activate
        pip install -r requirements.txt
        pip install bandit
        bandit -r app -x .venv,tests,__pycache__,**/site-packages/** -ll
        
    elif [ -f "go.mod" ]; then
        echo "Running gosec for Go..."
        go install github.com/securego/gosec/v2/cmd/gosec@latest
        export PATH=$(go env GOPATH)/bin:$PATH
        gosec -severity medium -confidence medium -fmt json -out gosec-report.json ./...
    else
        echo "Unknown service type for $SERVICE"
        exit 1
    fi
    
    cd ..
    echo "✅ Security scan completed for $SERVICE"
done