#!/bin/bash

# Exit script on error
set -e

# Function to deploy the Go application
deploy_application() {
    echo "Starting deployment..."

    # Pull the latest code from the Git repository
    echo "Pulling the latest code..."
    git pull origin main

    # Build the Go application
    echo "Building the Go application..."
    go build -o bin/app ./cmd/api

    # Check if Docker is used and rebuild/restart containers if necessary
    if [ -f "docker-compose.yml" ]; then
        echo "Docker Compose detected. Rebuilding containers..."
        docker-compose down
        docker-compose build
        docker-compose up -d
    else
        echo "Starting the application directly..."
        # Stop the running application (if any)
        if pgrep -f "bin/app" > /dev/null; then
            echo "Stopping running application..."
            pkill -f "bin/app"
        fi

        # Start the application
        echo "Starting the new version of the application..."
        nohup ./bin/app &> /dev/null &
    fi

    echo "Deployment completed!"
}

# Run deployment
deploy_application
