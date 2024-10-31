#!/usr/bin/env bash
# Create a .env file from Render environment variables

echo "Creating .env file from environment variables"

# Add each environment variable to the .env file
echo "PORT=${PORT}" >> .env
echo "MONGODB_URI=${MONGODB_URI}" >> .env
echo "ClientURL=${ClientURL}" >> .env
# Add as many as you need based on your environment variables

echo ".env file created successfully"
