#!/bin/bash

# Generate Swagger documentation
# Install swag: go install github.com/swaggo/swag/cmd/swag@latest

swag init -g cmd/server/main.go -o docs

