#!/bin/bash

# Setup script for RestySched

echo "=== RestySched Setup ==="
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.23 or higher."
    exit 1
fi

echo "✓ Go is installed: $(go version)"
echo ""

# Install Templ CLI
echo "Installing Templ CLI..."
go install github.com/a-h/templ/cmd/templ@latest
echo "✓ Templ CLI installed"
echo ""

# Download dependencies
echo "Downloading Go dependencies..."
go mod download
echo "✓ Dependencies downloaded"
echo ""

# Generate templates
echo "Generating Templ templates..."
templ generate
echo "✓ Templates generated"
echo ""

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo "Creating .env file from example..."
    cp .env.example .env
    echo "✓ .env file created"
    echo ""
    echo "⚠️  Please update .env file with your n8n webhook URL!"
else
    echo "✓ .env file already exists"
fi

echo ""
echo "=== Setup Complete! ==="
echo ""
echo "Next steps:"
echo "1. Update .env file with your n8n webhook URL"
echo "2. Run: make run"
echo "3. Access the app at http://localhost:8080"
echo ""
