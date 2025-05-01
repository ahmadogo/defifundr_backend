#!/bin/bash

# Navigate to project root directory
cd /Users/ademolakolawole/Projects/defifundr_backend

# Check if main.go exists in expected location
if [ -f "./cmd/api/main.go" ]; then
  echo "✅ Found main.go in expected location"
else
  echo "❌ main.go not found at ./cmd/api/main.go"
  echo "Searching for main.go files:"
  find . -name "main.go"
fi

# Check if waitlist.sql.go exists and for Swag annotations
if [ -f "./db/sqlc/waitlist.sql.go" ]; then
  echo "✅ Found waitlist.sql.go"
  
  # Check for relevant annotations
  if grep -q "@Param status" "./db/sqlc/waitlist.sql.go"; then
    echo "✅ Found @Param status annotation in waitlist.sql.go"
  else
    echo "❌ Missing @Param status annotation in waitlist.sql.go"
  fi
else
  echo "❌ waitlist.sql.go not found at ./db/sqlc/waitlist.sql.go"
fi

# Check Go module setup
if [ -f "go.mod" ]; then
  echo "✅ Found go.mod file"
  echo "Module name: $(grep "^module" go.mod)"
else
  echo "❌ go.mod not found. Is this a Go module?"
fi

# Check where Swag is looking for Go files
echo "Checking Go files in current directory:"
ls -la *.go 2>/dev/null || echo "No Go files in root directory"

# Check Swag version
echo "Swag version: $(swag --version)"