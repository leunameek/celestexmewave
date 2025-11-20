@echo off
REM CelestexMewave Backend Start Script (Windows)
REM This script starts the Go backend server with Supabase

setlocal enabledelayedexpansion

echo.
echo 🚀 Starting CelestexMewave Backend...
echo.

REM Check if .env file exists
if not exist "backend\.env" (
    echo ❌ Error: backend\.env file not found!
    echo.
    echo Please create backend\.env with the following content:
    echo.
    echo DATABASE_URL=postgresql://postgres:[YOUR_PASSWORD]@[YOUR_PROJECT_ID].supabase.co:5432/postgres?sslmode=require
    echo DB_DRIVER=postgres
    echo SERVER_PORT=8080
    echo SERVER_ENV=development
    echo SERVER_HOST=0.0.0.0
    echo JWT_SECRET=your_super_secret_key_change_this_in_production
    echo JWT_EXPIRATION=24h
    echo REFRESH_TOKEN_EXPIRATION=7d
    echo SMTP_HOST=smtp.gmail.com
    echo SMTP_PORT=587
    echo SMTP_USER=your_email@gmail.com
    echo SMTP_PASSWORD=your_app_specific_password
    echo SMTP_FROM=noreply@celestexmewave.com
    echo FRONTEND_URL=http://localhost:3000
    echo UPLOAD_DIR=./uploads/images
    echo MAX_UPLOAD_SIZE=5242880
    echo.
    pause
    exit /b 1
)

REM Check if Go is installed
go version >nul 2>&1
if errorlevel 1 (
    echo ❌ Error: Go is not installed!
    echo Please install Go 1.25.3 or higher from https://golang.org/dl/
    echo.
    pause
    exit /b 1
)

echo ✓ Configuration found
echo ✓ Go is installed
echo.

REM Navigate to backend directory
cd backend

REM Download dependencies
echo 📦 Downloading dependencies...
go mod download
go mod tidy
echo ✓ Dependencies ready
echo.

REM Start the server
echo 🔧 Starting server...
echo.
go run cmd/api/main.go

pause
